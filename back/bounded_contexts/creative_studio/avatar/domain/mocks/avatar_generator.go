package mocks

import (
	avatardomain "brandtoonapi/bounded_contexts/creative_studio/avatar/domain"
	"context"
)

type AvatarGeneratorMock struct {
	GenerateOptionsFunc func(ctx context.Context, prompt string, count int) ([]avatardomain.GeneratedAvatarImage, error)
	GenerateCalls       int
	LastPrompt          string
	LastCount           int
}

func (m *AvatarGeneratorMock) GenerateOptions(
	ctx context.Context,
	prompt string,
	count int,
) ([]avatardomain.GeneratedAvatarImage, error) {
	m.GenerateCalls++
	m.LastPrompt = prompt
	m.LastCount = count
	if m.GenerateOptionsFunc == nil {
		return nil, nil
	}

	return m.GenerateOptionsFunc(ctx, prompt, count)
}
