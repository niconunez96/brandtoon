package userrepo

import (
	userdomain "brandtoonapi/bounded_contexts/identity/user/domain"
	"time"
)

type userDBModel struct {
	AvatarURL string     `db:"avatar_url"`
	CreatedAt time.Time  `db:"created_at"`
	DeletedAt *time.Time `db:"deleted_at"`
	Email     string     `db:"email"`
	ID        string     `db:"id"`
	Name      string     `db:"name"`
	UpdatedAt time.Time  `db:"updated_at"`
}

func newUserDBModel(user userdomain.User) *userDBModel {
	return &userDBModel{
		AvatarURL: user.AvatarURL,
		Email:     user.Email,
		ID:        user.ID,
		Name:      user.Name,
	}
}

func (m *userDBModel) GetID() string {
	return m.ID
}

func (m *userDBModel) InsertValues() map[string]any {
	return map[string]any{
		"avatar_url": m.AvatarURL,
		"created_at": m.CreatedAt,
		"deleted_at": m.DeletedAt,
		"email":      m.Email,
		"id":         m.ID,
		"name":       m.Name,
		"updated_at": m.UpdatedAt,
	}
}

func (m *userDBModel) SetCreatedAt(createdAt time.Time) {
	m.CreatedAt = createdAt
}

func (m *userDBModel) SetDeletedAt(deletedAt *time.Time) {
	m.DeletedAt = deletedAt
}

func (m *userDBModel) SetUpdatedAt(updatedAt time.Time) {
	m.UpdatedAt = updatedAt
}

func (m *userDBModel) TableName() string {
	return "users"
}

func (m *userDBModel) ToDomain() userdomain.User {
	return userdomain.NewUser(m.ID, m.Email, m.Name, m.AvatarURL)
}

func (m *userDBModel) UpdateValues() map[string]any {
	return map[string]any{
		"avatar_url": m.AvatarURL,
		"email":      m.Email,
		"name":       m.Name,
		"updated_at": m.UpdatedAt,
	}
}
