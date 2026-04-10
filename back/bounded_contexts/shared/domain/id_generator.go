package domain

import "github.com/google/uuid"

type IDGenerator func() (string, error)

func GenerateUUIDv7() (string, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return "", err
	}

	return id.String(), nil
}
