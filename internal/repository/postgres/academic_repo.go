package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"mentorhub/internal/domain"
	"mentorhub/internal/repository"
)

// ── Academic Years ─────────────────────────────────────────────────────────────

type academicYearRepository struct {
	db *pgxpool.Pool
}

func NewAcademicYearRepository(db *pgxpool.Pool) repository.AcademicYearRepository {
	return &academicYearRepository{db: db}
}

func (r *academicYearRepository) Create(ctx context.Context, ay *domain.AcademicYear) error {
	const q = `
		INSERT INTO academic_years (id, name, start_date, end_date, is_active)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING created_at
	`
	if ay.ID == uuid.Nil {
		ay.ID = uuid.New()
	}
	return r.db.QueryRow(ctx, q,
		ay.ID, ay.Name, ay.StartDate, ay.EndDate, ay.IsActive,
	).Scan(&ay.CreatedAt)
}

func (r *academicYearRepository) GetAll(ctx context.Context) ([]*domain.AcademicYear, error) {
	const q = `
		SELECT id, name, start_date, end_date, is_active, created_at
		FROM academic_years ORDER BY start_date DESC
	`
	rows, err := r.db.Query(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("query academic years: %w", err)
	}
	defer rows.Close()

	var years []*domain.AcademicYear
	for rows.Next() {
		var ay domain.AcademicYear
		if err := rows.Scan(&ay.ID, &ay.Name, &ay.StartDate, &ay.EndDate, &ay.IsActive, &ay.CreatedAt); err != nil {
			return nil, err
		}
		years = append(years, &ay)
	}
	return years, rows.Err()
}

func (r *academicYearRepository) GetActive(ctx context.Context) (*domain.AcademicYear, error) {
	const q = `
		SELECT id, name, start_date, end_date, is_active, created_at
		FROM academic_years WHERE is_active = TRUE LIMIT 1
	`
	var ay domain.AcademicYear
	err := r.db.QueryRow(ctx, q).Scan(&ay.ID, &ay.Name, &ay.StartDate, &ay.EndDate, &ay.IsActive, &ay.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("no active academic year")
		}
		return nil, fmt.Errorf("get active academic year: %w", err)
	}
	return &ay, nil
}

func (r *academicYearRepository) SetActive(ctx context.Context, id uuid.UUID) error {
	// Транзакция: снять флаг со всех → поставить на один
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, `UPDATE academic_years SET is_active = FALSE`); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, `UPDATE academic_years SET is_active = TRUE WHERE id = $1`, id); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// ── Faculties & Specialties ───────────────────────────────────────────────────

type facultyRepository struct {
	db *pgxpool.Pool
}

func NewFacultyRepository(db *pgxpool.Pool) repository.FacultyRepository {
	return &facultyRepository{db: db}
}

func (r *facultyRepository) CreateFaculty(ctx context.Context, f *domain.Faculty) error {
	const q = `
		INSERT INTO faculties (id, name, code)
		VALUES ($1, $2, $3)
		RETURNING created_at
	`
	if f.ID == uuid.Nil {
		f.ID = uuid.New()
	}
	return r.db.QueryRow(ctx, q, f.ID, f.Name, f.Code).Scan(&f.CreatedAt)
}

func (r *facultyRepository) GetAllFaculties(ctx context.Context) ([]*domain.Faculty, error) {
	const q = `
		SELECT f.id, f.name, f.code, f.created_at,
		       s.id, s.name, s.created_at
		FROM faculties f
		LEFT JOIN specialties s ON s.faculty_id = f.id
		ORDER BY f.name, s.name
	`
	rows, err := r.db.Query(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("query faculties: %w", err)
	}
	defer rows.Close()

	facultyMap := make(map[uuid.UUID]*domain.Faculty)
	var order []uuid.UUID

	for rows.Next() {
		var (
			f    domain.Faculty
			sID  *uuid.UUID
			sNm  *string
			sCAt *time.Time
		)
		if err := rows.Scan(&f.ID, &f.Name, &f.Code, &f.CreatedAt, &sID, &sNm, &sCAt); err != nil {
			return nil, err
		}
		if _, ok := facultyMap[f.ID]; !ok {
			fp := f
			fp.Specialties = []*domain.Specialty{}
			facultyMap[f.ID] = &fp
			order = append(order, f.ID)
		}
		if sID != nil {
			facultyMap[f.ID].Specialties = append(facultyMap[f.ID].Specialties, &domain.Specialty{
				ID:        *sID,
				FacultyID: f.ID,
				Name:      *sNm,
				CreatedAt: *sCAt,
			})
		}
	}

	result := make([]*domain.Faculty, 0, len(order))
	for _, id := range order {
		result = append(result, facultyMap[id])
	}
	return result, rows.Err()
}

func (r *facultyRepository) CreateSpecialty(ctx context.Context, s *domain.Specialty) error {
	const q = `
		INSERT INTO specialties (id, faculty_id, name)
		VALUES ($1, $2, $3)
		RETURNING created_at
	`
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return r.db.QueryRow(ctx, q, s.ID, s.FacultyID, s.Name).Scan(&s.CreatedAt)
}

func (r *facultyRepository) GetSpecialtiesByFaculty(ctx context.Context, facultyID uuid.UUID) ([]*domain.Specialty, error) {
	const q = `
		SELECT id, faculty_id, name, created_at
		FROM specialties WHERE faculty_id = $1 ORDER BY name
	`
	rows, err := r.db.Query(ctx, q, facultyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var specs []*domain.Specialty
	for rows.Next() {
		var s domain.Specialty
		if err := rows.Scan(&s.ID, &s.FacultyID, &s.Name, &s.CreatedAt); err != nil {
			return nil, err
		}
		specs = append(specs, &s)
	}
	return specs, rows.Err()
}
