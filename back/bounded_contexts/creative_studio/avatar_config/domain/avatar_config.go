package avatarconfigdomain

import "errors"

var ErrInvalidArtisticStyle = errors.New("invalid artistic style")

type ArtisticStyle string

const (
	ArtisticStyle2D ArtisticStyle = "2D"
	ArtisticStyle3D ArtisticStyle = "3D"
)

type AvatarConfig struct {
	AvatarID      string
	ArtisticStyle ArtisticStyle
	Prompt        string
}

func NewAvatarConfig(
	avatarID string,
	prompt string,
	artisticStyle ArtisticStyle,
) AvatarConfig {
	return AvatarConfig{
		AvatarID:      avatarID,
		ArtisticStyle: artisticStyle,
		Prompt:        prompt,
	}
}

func ParseArtisticStyle(value string) (ArtisticStyle, error) {
	switch ArtisticStyle(value) {
	case ArtisticStyle2D:
		return ArtisticStyle2D, nil
	case ArtisticStyle3D:
		return ArtisticStyle3D, nil
	default:
		return "", ErrInvalidArtisticStyle
	}
}
