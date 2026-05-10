package test

import (
	"testing"

	avatardomain "brandtoonapi/bounded_contexts/creative_studio/avatar/domain"
	avatardto "brandtoonapi/bounded_contexts/creative_studio/avatar/useCases/dto"
)

func TestSerializeReturnsIDAndName(t *testing.T) {
	t.Parallel()

	serialized := avatardto.Serialize(avatardomain.NewAvatar("avatar-v7", "user-v7", "Studio Hero"))

	if serialized.ID != "avatar-v7" || serialized.Name != "Studio Hero" {
		t.Fatalf("expected serialized avatar identity, got %+v", serialized)
	}
}

func TestSerializeListPreservesOrderAndSupportsEmptySlices(t *testing.T) {
	t.Parallel()

	t.Run("preserves repository order", func(t *testing.T) {
		avatars := []avatardomain.Avatar{
			avatardomain.NewAvatar("avatar-2", "user-v7", "Beta"),
			avatardomain.NewAvatar("avatar-1", "user-v7", "Alpha"),
		}

		serialized := avatardto.SerializeList(avatars)

		if len(serialized) != 2 || serialized[0].ID != "avatar-2" || serialized[1].ID != "avatar-1" {
			t.Fatalf("expected serialized avatars to preserve order, got %+v", serialized)
		}
	})

	t.Run("returns empty slice for empty input", func(t *testing.T) {
		serialized := avatardto.SerializeList(nil)

		if serialized == nil || len(serialized) != 0 {
			t.Fatalf("expected non-nil empty slice, got %+v", serialized)
		}
	})
}

func TestSerializeReusesCanonicalAvatarDTOForDetailsReads(t *testing.T) {
	t.Parallel()

	avatar := avatardomain.NewAvatar("avatar-v7", "user-v7", "Studio Hero")

	serialized := avatardto.Serialize(avatar)

	if serialized.ID != "avatar-v7" || serialized.Name != "Studio Hero" {
		t.Fatalf("expected serialized avatar details identity, got %+v", serialized)
	}
}
