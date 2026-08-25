package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"mentorhub/internal/domain"
)

// ── Meetings ──────────────────────────────────────────────────────────────────

// MeetingRepository управляет встречами и объявлениями ментора.
type MeetingRepository struct {
	db *pgxpool.Pool
}

func NewMeetingRepository(db *pgxpool.Pool) *MeetingRepository {
	return &MeetingRepository{db: db}
}

func (r *MeetingRepository) CreateMeeting(ctx context.Context, m *domain.Meeting) error {
	const q = `
		INSERT INTO meetings (id, mentor_id, group_id, title, description, scheduled_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING created_at
	`
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return r.db.QueryRow(ctx, q,
		m.ID, m.MentorID, m.GroupID, m.Title, m.Description, m.ScheduledAt,
	).Scan(&m.CreatedAt)
}

func (r *MeetingRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Meeting, error) {
	const q = `
		SELECT id, mentor_id, group_id, title, description, scheduled_at, held, notes, created_at
		FROM meetings WHERE id = $1
	`
	var m domain.Meeting
	err := r.db.QueryRow(ctx, q, id).Scan(
		&m.ID, &m.MentorID, &m.GroupID, &m.Title, &m.Description,
		&m.ScheduledAt, &m.Held, &m.Notes, &m.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("meeting %s not found", id)
		}
		return nil, err
	}
	return &m, nil
}

// GetMentorMeetings возвращает все встречи ментора (ближайшие вверху).
func (r *MeetingRepository) GetMentorMeetings(ctx context.Context, mentorID uuid.UUID) ([]*domain.Meeting, error) {
	const q = `
		SELECT id, mentor_id, group_id, title, description, scheduled_at, held, notes, created_at
		FROM meetings
		WHERE mentor_id = $1
		ORDER BY scheduled_at DESC
	`
	rows, err := r.db.Query(ctx, q, mentorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var meetings []*domain.Meeting
	for rows.Next() {
		var m domain.Meeting
		if err := rows.Scan(
			&m.ID, &m.MentorID, &m.GroupID, &m.Title, &m.Description,
			&m.ScheduledAt, &m.Held, &m.Notes, &m.CreatedAt,
		); err != nil {
			return nil, err
		}
		meetings = append(meetings, &m)
	}
	return meetings, rows.Err()
}

// CompleteMeeting помечает встречу как проведённую и добавляет заметки.
func (r *MeetingRepository) CompleteMeeting(ctx context.Context, meetingID uuid.UUID, mentorID uuid.UUID, notes string) error {
	const q = `
		UPDATE meetings
		SET held = TRUE, notes = $3
		WHERE id = $1 AND mentor_id = $2 AND held = FALSE
	`
	tag, err := r.db.Exec(ctx, q, meetingID, mentorID, notes)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("meeting not found, already completed, or you are not the owner")
	}
	return nil
}

// ── Announcements ─────────────────────────────────────────────────────────────

func (r *MeetingRepository) CreateAnnouncement(ctx context.Context, a *domain.Announcement) error {
	const q = `
		INSERT INTO announcements (id, author_id, group_id, title, body)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING created_at
	`
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return r.db.QueryRow(ctx, q,
		a.ID, a.AuthorID, a.GroupID, a.Title, a.Body,
	).Scan(&a.CreatedAt)
}

// GetGroupAnnouncements возвращает объявления группы + глобальные (group_id IS NULL).
func (r *MeetingRepository) GetGroupAnnouncements(ctx context.Context, groupID *uuid.UUID) ([]*domain.Announcement, error) {
	const q = `
		SELECT a.id, a.author_id, a.group_id, a.title, a.body, a.created_at,
		       u.first_name, u.last_name
		FROM announcements a
		JOIN users u ON u.id = a.author_id
		WHERE a.group_id = $1 OR a.group_id IS NULL
		ORDER BY a.created_at DESC
	`
	rows, err := r.db.Query(ctx, q, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var announcements []*domain.Announcement
	for rows.Next() {
		var a domain.Announcement
		var firstName, lastName string
		if err := rows.Scan(
			&a.ID, &a.AuthorID, &a.GroupID, &a.Title, &a.Body, &a.CreatedAt,
			&firstName, &lastName,
		); err != nil {
			return nil, err
		}
		a.Author = &domain.User{ID: a.AuthorID, FirstName: firstName, LastName: lastName}
		announcements = append(announcements, &a)
	}
	return announcements, rows.Err()
}

// ── Notifications ─────────────────────────────────────────────────────────────

// NotificationRepository реализует repository.NotificationRepository через pgxpool.
type NotificationRepository struct {
	db *pgxpool.Pool
}

func NewNotificationRepository(db *pgxpool.Pool) *NotificationRepository {
	return &NotificationRepository{db: db}
}

func (r *NotificationRepository) Create(ctx context.Context, n *domain.Notification) error {
	const q = `
		INSERT INTO notifications (id, user_id, title, body, is_read)
		VALUES ($1, $2, $3, $4, FALSE)
		RETURNING created_at
	`
	if n.ID == uuid.Nil {
		n.ID = uuid.New()
	}
	return r.db.QueryRow(ctx, q, n.ID, n.UserID, n.Title, n.Body).Scan(&n.CreatedAt)
}

func (r *NotificationRepository) GetByUserID(ctx context.Context, userID uuid.UUID, unreadOnly bool) ([]*domain.Notification, error) {
	q := `
		SELECT id, user_id, title, body, is_read, created_at
		FROM notifications
		WHERE user_id = $1
	`
	if unreadOnly {
		q += ` AND is_read = FALSE`
	}
	q += ` ORDER BY created_at DESC LIMIT 100`

	rows, err := r.db.Query(ctx, q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []*domain.Notification
	for rows.Next() {
		var n domain.Notification
		if err := rows.Scan(&n.ID, &n.UserID, &n.Title, &n.Body, &n.IsRead, &n.CreatedAt); err != nil {
			return nil, err
		}
		notifications = append(notifications, &n)
	}
	return notifications, rows.Err()
}

func (r *NotificationRepository) MarkAsRead(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, `UPDATE notifications SET is_read = TRUE WHERE id = $1`, id)
	return err
}

func (r *NotificationRepository) MarkAllAsRead(ctx context.Context, userID uuid.UUID) error {
	_, err := r.db.Exec(ctx,
		`UPDATE notifications SET is_read = TRUE WHERE user_id = $1 AND is_read = FALSE`,
		userID,
	)
	return err
}

func (r *NotificationRepository) CountUnread(ctx context.Context, userID uuid.UUID) (int, error) {
	var count int
	err := r.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM notifications WHERE user_id = $1 AND is_read = FALSE`,
		userID,
	).Scan(&count)
	return count, err
}
