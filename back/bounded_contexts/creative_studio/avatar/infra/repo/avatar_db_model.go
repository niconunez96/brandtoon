package avatarrepo

import (
	avatardomain "brandtoonapi/bounded_contexts/creative_studio/avatar/domain"
	"time"
)

type avatarDBModel struct {
	CreatedAt time.Time  `db:"created_at"`
	DeletedAt *time.Time `db:"deleted_at"`
	ID        string     `db:"id"`
	Name      string     `db:"name"`
	UpdatedAt time.Time  `db:"updated_at"`
	UserID    string     `db:"user_id"`
}

func newAvatarDBModel(avatar avatardomain.Avatar) *avatarDBModel {
	return &avatarDBModel{
		ID:     avatar.ID,
		Name:   avatar.Name,
		UserID: avatar.UserID,
	}
}

func (m *avatarDBModel) GetID() string {
	return m.ID
}

func (m *avatarDBModel) InsertValues() map[string]any {
	return map[string]any{
		"created_at": m.CreatedAt,
		"deleted_at": m.DeletedAt,
		"id":         m.ID,
		"name":       m.Name,
		"updated_at": m.UpdatedAt,
		"user_id":    m.UserID,
	}
}

func (m *avatarDBModel) SetCreatedAt(createdAt time.Time) {
	m.CreatedAt = createdAt
}

func (m *avatarDBModel) SetDeletedAt(deletedAt *time.Time) {
	m.DeletedAt = deletedAt
}

func (m *avatarDBModel) SetUpdatedAt(updatedAt time.Time) {
	m.UpdatedAt = updatedAt
}

func (m *avatarDBModel) TableName() string {
	return "avatars"
}

func (m *avatarDBModel) ToDomain() avatardomain.Avatar {
	return avatardomain.NewAvatar(m.ID, m.UserID, m.Name)
}

func (m *avatarDBModel) UpdateValues() map[string]any {
	return map[string]any{
		"name":       m.Name,
		"updated_at": m.UpdatedAt,
	}
}
