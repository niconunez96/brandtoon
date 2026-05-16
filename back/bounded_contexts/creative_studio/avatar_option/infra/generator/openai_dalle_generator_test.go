package avataroptiongenerator

import (
	"context"
	"testing"

	"github.com/openai/openai-go/v3"
)

type imageGeneratorClientStub struct {
	called     bool
	lastParams openai.ImageGenerateParams
}

func (s *imageGeneratorClientStub) Generate(
	ctx context.Context,
	params openai.ImageGenerateParams,
) (*openai.ImagesResponse, error) {
	s.called = true
	s.lastParams = params
	return &openai.ImagesResponse{}, nil
}

func TestGenerateOptionsReturnsEmptySliceWithoutCallingClientWhenCountIsZero(t *testing.T) {
	t.Parallel()

	client := &imageGeneratorClientStub{}
	generator := &OpenAIDALLEGenerator{client: client}

	images, err := generator.GenerateOptions(context.Background(), "hero prompt", 0)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if client.called {
		t.Fatal("expected client not to be called for zero count")
	}
	if images == nil || len(images) != 0 {
		t.Fatalf("expected non-nil empty slice, got %+v", images)
	}
}
