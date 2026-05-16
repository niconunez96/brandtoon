package avataroptiondomain

import "context"

type GeneratedAvatarImage struct {
	ContentType string
	Data        []byte
}

type AvatarGenerator interface {
	GenerateOptions(ctx context.Context, prompt string, count int) ([]GeneratedAvatarImage, error)
}
