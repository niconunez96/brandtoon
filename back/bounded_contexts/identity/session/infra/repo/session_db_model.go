package sessionrepo

import (
	"time"

	"brandtoonapi/bounded_contexts/identity/session/domain"
)

type sessionDBModel struct {
	CreatedAt time.Time  `db:"created_at"`
	DeletedAt *time.Time `db:"deleted_at"`
	ExpiresAt time.Time  `db:"expires_at"`
	ID        string     `db:"id"`
	UpdatedAt time.Time  `db:"updated_at"`
	UserID    string     `db:"user_id"`
}

func newSessionDBModel(session sessiondomain.Session) *sessionDBModel {
	return &sessionDBModel{
		ExpiresAt: session.ExpiresAt,
		ID:        session.ID,
		UserID:    session.UserID,
	}
}

func (m *sessionDBModel) GetID() string {
	return m.ID
}

func (m *sessionDBModel) InsertValues() map[string]any {
	return map[string]any{
		"created_at": m.CreatedAt,
		"deleted_at": m.DeletedAt,
		"expires_at": m.ExpiresAt,
		"id":         m.ID,
		"updated_at": m.UpdatedAt,
		"user_id":    m.UserID,
	}
}

func (m *sessionDBModel) SetCreatedAt(createdAt time.Time) {
	m.CreatedAt = createdAt
}

func (m *sessionDBModel) SetDeletedAt(deletedAt *time.Time) {
	m.DeletedAt = deletedAt
}

func (m *sessionDBModel) SetUpdatedAt(updatedAt time.Time) {
	m.UpdatedAt = updatedAt
}

func (m *sessionDBModel) TableName() string {
	return "sessions"
}

func (m *sessionDBModel) ToDomain() sessiondomain.Session {
	return sessiondomain.NewSession(m.ID, m.UserID, m.ExpiresAt)
}

func (m *sessionDBModel) UpdateValues() map[string]any {
	return map[string]any{
		"expires_at": m.ExpiresAt,
		"updated_at": m.UpdatedAt,
		"user_id":    m.UserID,
	}
}
