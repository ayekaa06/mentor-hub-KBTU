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

// ── FAQ ───────────────────────────────────────────────────────────────────────

type faqRepository struct {
	db *pgxpool.Pool
}

// NewFAQRepository создаёт реализацию repository.FAQRepository.
func NewFAQRepository(db *pgxpool.Pool) repository.FAQRepository {
	return &faqRepository{db: db}
}

// GetActive возвращает активные FAQ-записи, отсортированные по order_num.
func (r *faqRepository) GetActive(ctx context.Context) ([]*domain.FAQItem, error) {
	const q = `
		SELECT id, question, answer, order_num, is_active, created_by, created_at, updated_at
		FROM faq_items
		WHERE is_active = TRUE
		ORDER BY order_num ASC
	`
	return r.scanFAQItems(ctx, q)
}

// GetAll возвращает все FAQ-записи (включая неактивные).
func (r *faqRepository) GetAll(ctx context.Context) ([]*domain.FAQItem, error) {
	const q = `
		SELECT id, question, answer, order_num, is_active, created_by, created_at, updated_at
		FROM faq_items
		ORDER BY order_num ASC
	`
	return r.scanFAQItems(ctx, q)
}

func (r *faqRepository) scanFAQItems(ctx context.Context, q string, args ...any) ([]*domain.FAQItem, error) {
	rows, err := r.db.Query(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("query faq items: %w", err)
	}
	defer rows.Close()

	var items []*domain.FAQItem
	for rows.Next() {
		var item domain.FAQItem
		if err := rows.Scan(
			&item.ID, &item.Question, &item.Answer, &item.OrderNum,
			&item.IsActive, &item.CreatedBy, &item.CreatedAt, &item.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	return items, rows.Err()
}

// ── Questions ─────────────────────────────────────────────────────────────────

type questionRepository struct {
	db *pgxpool.Pool
}

// NewQuestionRepository создаёт реализацию repository.QuestionRepository.
func NewQuestionRepository(db *pgxpool.Pool) repository.QuestionRepository {
	return &questionRepository{db: db}
}

// Create сохраняет новый вопрос freshman-а.
func (r *questionRepository) Create(ctx context.Context, q *domain.Question) error {
	const stmt = `
		INSERT INTO questions (id, freshman_id, mentor_id, body)
		VALUES ($1, $2, $3, $4)
		RETURNING created_at
	`
	if q.ID == uuid.Nil {
		q.ID = uuid.New()
	}
	return r.db.QueryRow(ctx, stmt,
		q.ID, q.FreshmanID, q.MentorID, q.Body,
	).Scan(&q.CreatedAt)
}

// GetByFreshmanID возвращает все вопросы freshman-а (с именем ментора).
func (r *questionRepository) GetByFreshmanID(ctx context.Context, freshmanID uuid.UUID) ([]*domain.Question, error) {
	const q = `
		SELECT qu.id, qu.freshman_id, qu.mentor_id, qu.body,
		       qu.answer, qu.answered_at, qu.created_at,
		       u.first_name, u.last_name
		FROM questions qu
		JOIN users u ON u.id = qu.mentor_id
		WHERE qu.freshman_id = $1
		ORDER BY qu.created_at DESC
	`
	rows, err := r.db.Query(ctx, q, freshmanID)
	if err != nil {
		return nil, fmt.Errorf("query questions by freshman: %w", err)
	}
	defer rows.Close()

	var questions []*domain.Question
	for rows.Next() {
		var qu domain.Question
		var mFirst, mLast string
		if err := rows.Scan(
			&qu.ID, &qu.FreshmanID, &qu.MentorID, &qu.Body,
			&qu.Answer, &qu.AnsweredAt, &qu.CreatedAt,
			&mFirst, &mLast,
		); err != nil {
			return nil, err
		}
		_ = mFirst
		_ = mLast
		questions = append(questions, &qu)
	}
	return questions, rows.Err()
}

// GetByMentorID возвращает все вопросы, адресованные ментору (с именем freshman-а).
func (r *questionRepository) GetByMentorID(ctx context.Context, mentorID uuid.UUID) ([]*domain.Question, error) {
	const q = `
		SELECT qu.id, qu.freshman_id, qu.mentor_id, qu.body,
		       qu.answer, qu.answered_at, qu.created_at,
		       u.first_name, u.last_name
		FROM questions qu
		JOIN users u ON u.id = qu.freshman_id
		WHERE qu.mentor_id = $1
		ORDER BY qu.created_at DESC
	`
	rows, err := r.db.Query(ctx, q, mentorID)
	if err != nil {
		return nil, fmt.Errorf("query questions by mentor: %w", err)
	}
	defer rows.Close()

	var questions []*domain.Question
	for rows.Next() {
		var qu domain.Question
		var firstName, lastName string
		if err := rows.Scan(
			&qu.ID, &qu.FreshmanID, &qu.MentorID, &qu.Body,
			&qu.Answer, &qu.AnsweredAt, &qu.CreatedAt,
			&firstName, &lastName,
		); err != nil {
			return nil, err
		}
		qu.Freshman = &domain.User{ID: qu.FreshmanID, FirstName: firstName, LastName: lastName}
		questions = append(questions, &qu)
	}
	return questions, rows.Err()
}

// Answer записывает ответ ментора на вопрос freshman-а.
func (r *questionRepository) Answer(ctx context.Context, questionID uuid.UUID, mentorID uuid.UUID, answer string) error {
	now := time.Now()
	const q = `
		UPDATE questions
		SET answer = $3, answered_at = $4
		WHERE id = $1 AND mentor_id = $2 AND answer IS NULL
	`
	tag, err := r.db.Exec(ctx, q, questionID, mentorID, answer, now)
	if err != nil {
		return fmt.Errorf("answer question: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("question not found, already answered, or you are not the mentor")
	}
	return nil
}

// GetByID возвращает вопрос по ID (вспомогательный метод).
func (r *questionRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Question, error) {
	const q = `
		SELECT id, freshman_id, mentor_id, body, answer, answered_at, created_at
		FROM questions WHERE id = $1
	`
	var qu domain.Question
	err := r.db.QueryRow(ctx, q, id).Scan(
		&qu.ID, &qu.FreshmanID, &qu.MentorID, &qu.Body,
		&qu.Answer, &qu.AnsweredAt, &qu.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("question %s not found", id)
		}
		return nil, err
	}
	return &qu, nil
}
