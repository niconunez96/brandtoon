package avatargenerator

import (
	avatardomain "brandtoonapi/bounded_contexts/creative_studio/avatar/domain"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

const defaultOpenAIImageModel = "dall-e-3"

type imageGeneratorClient interface {
	Generate(ctx context.Context, params openai.ImageGenerateParams) (*openai.ImagesResponse, error)
}

type openAIImagesClient struct {
	client *openai.Client
}

func (c openAIImagesClient) Generate(
	ctx context.Context,
	params openai.ImageGenerateParams,
) (*openai.ImagesResponse, error) {
	return c.client.Images.Generate(ctx, params)
}

type OpenAIDALLEGenerator struct {
	client imageGeneratorClient
	model  string
}

func NewOpenAIDALLEGenerator(apiKey string, model string) (*OpenAIDALLEGenerator, error) {
	if strings.TrimSpace(apiKey) == "" {
		return nil, errors.New("openai api key is required")
	}

	trimmedModel := strings.TrimSpace(model)
	if trimmedModel == "" {
		trimmedModel = defaultOpenAIImageModel
	}

	client := openai.NewClient(option.WithAPIKey(apiKey))
	return &OpenAIDALLEGenerator{
		client: openAIImagesClient{client: &client},
		model:  trimmedModel,
	}, nil
}

func (g *OpenAIDALLEGenerator) GenerateOptions(
	ctx context.Context,
	prompt string,
	count int,
) ([]avatardomain.GeneratedAvatarImage, error) {
	if count <= 0 {
		return []avatardomain.GeneratedAvatarImage{}, nil
	}

	response, err := g.client.Generate(ctx, openai.ImageGenerateParams{
		Prompt:         prompt,
		Model:          openai.ImageModel(g.model),
		N:              openai.Int(int64(count)),
		// ResponseFormat: openai.ImageGenerateParamsResponseFormatB64JSON,
		Size:           openai.ImageGenerateParamsSize1024x1024,
	})
	if err != nil {
		return nil, err
	}

	images := make([]avatardomain.GeneratedAvatarImage, 0, len(response.Data))
	for _, item := range response.Data {
		if item.B64JSON == "" {
			return nil, errors.New("openai image response missing b64 payload")
		}

		decoded, decodeErr := base64.StdEncoding.DecodeString(item.B64JSON)
		if decodeErr != nil {
			return nil, fmt.Errorf("decode openai image: %w", decodeErr)
		}

		images = append(images, avatardomain.GeneratedAvatarImage{
			ContentType: "image/png",
			Data:        decoded,
		})
	}

	return images, nil
}
