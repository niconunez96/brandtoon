package mocks

import (
	avataroptiondomain "brandtoonapi/bounded_contexts/creative_studio/avatar_option/domain"
	"context"
)

type AvatarGeneratorMock struct {
	GenerateOptionsFunc func(ctx context.Context, prompt string, count int) ([]avataroptiondomain.GeneratedAvatarImage, error)
	GenerateCalls       int
	LastPrompt          string
	LastCount           int
}

func (m *AvatarGeneratorMock) GenerateOptions(
	ctx context.Context,
	prompt string,
	count int,
) ([]avataroptiondomain.GeneratedAvatarImage, error) {
	m.GenerateCalls++
	m.LastPrompt = prompt
	m.LastCount = count
	if m.GenerateOptionsFunc == nil {
		return nil, nil
	}

	return m.GenerateOptionsFunc(ctx, prompt, count)
}
