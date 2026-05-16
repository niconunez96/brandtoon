package avatardomain

import "errors"

var ErrInvalidName = errors.New("invalid avatar name")

type Avatar struct {
	ID     string
	Name   string
	UserID string
}

func NewAvatar(id string, userID string, name string) Avatar {
	return Avatar{
		ID:     id,
		Name:   name,
		UserID: userID,
	}
}
