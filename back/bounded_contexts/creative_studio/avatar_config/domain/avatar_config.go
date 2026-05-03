package avatarconfigdomain

import "errors"

var ErrInvalidArtisticStyle = errors.New("invalid artistic style")
var ErrInvalidPersonality = errors.New("invalid personality")

type ArtisticStyle string

type Personality string

const (
	ArtisticStyle2D ArtisticStyle = "2D"
	ArtisticStyle3D ArtisticStyle = "3D"

	PersonalityFriendly Personality = "Friendly"
	PersonalityBold     Personality = "Bold"
	PersonalityPlayful  Personality = "Playful"
)

type AvatarConfig struct {
	AvatarID      string
	ArtisticStyle ArtisticStyle
	Personality   Personality
	Prompt        string
}

func NewAvatarConfig(
	avatarID string,
	prompt string,
	artisticStyle ArtisticStyle,
	personality Personality,
) AvatarConfig {
	return AvatarConfig{
		AvatarID:      avatarID,
		ArtisticStyle: artisticStyle,
		Personality:   personality,
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

func ParsePersonality(value string) (Personality, error) {
	switch Personality(value) {
	case PersonalityFriendly:
		return PersonalityFriendly, nil
	case PersonalityBold:
		return PersonalityBold, nil
	case PersonalityPlayful:
		return PersonalityPlayful, nil
	default:
		return "", ErrInvalidPersonality
	}
}
