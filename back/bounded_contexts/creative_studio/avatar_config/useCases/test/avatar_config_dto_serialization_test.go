package test

import (
	"testing"

	avatarconfigdomain "brandtoonapi/bounded_contexts/creative_studio/avatar_config/domain"
	avatarconfigdto "brandtoonapi/bounded_contexts/creative_studio/avatar_config/useCases/dto"
)

func TestSerializeMapsConfigFieldsAndNormalizesAvatarOptions(t *testing.T) {
	t.Parallel()

	avatarConfig := avatarconfigdomain.NewAvatarConfig(
		"avatar-v7",
		"Energetic mascot with bold shapes",
		avatarconfigdomain.ArtisticStyle3D,
		avatarconfigdomain.PersonalityBold,
	)

	serialized := avatarconfigdto.Serialize(avatarConfig)

	if serialized.AvatarID != "avatar-v7" ||
		serialized.ArtisticStyle != string(avatarconfigdomain.ArtisticStyle3D) {
		t.Fatalf("expected serialized config identity and artistic style, got %+v", serialized)
	}

	if serialized.Personality != string(avatarconfigdomain.PersonalityBold) ||
		serialized.Prompt != "Energetic mascot with bold shapes" {
		t.Fatalf("expected serialized config payload, got %+v", serialized)
	}
}
