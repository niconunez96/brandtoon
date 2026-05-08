package test

import (
	"testing"

	avatardomain "brandtoonapi/bounded_contexts/creative_studio/avatar/domain"
	avatardto "brandtoonapi/bounded_contexts/creative_studio/avatar/useCases/dto"
)

func TestSerializeAvatarReturnsIDAndName(t *testing.T) {
	t.Parallel()

	serialized := avatardto.SerializeAvatar(avatardomain.NewAvatar("avatar-v7", "user-v7", "Studio Hero"))

	if serialized.ID != "avatar-v7" || serialized.Name != "Studio Hero" {
		t.Fatalf("expected serialized avatar identity, got %+v", serialized)
	}
}

func TestSerializeAvatarListPreservesOrderAndSupportsEmptySlices(t *testing.T) {
	t.Parallel()

	t.Run("preserves repository order", func(t *testing.T) {
		avatars := []avatardomain.Avatar{
			avatardomain.NewAvatar("avatar-2", "user-v7", "Beta"),
			avatardomain.NewAvatar("avatar-1", "user-v7", "Alpha"),
		}

		serialized := avatardto.SerializeAvatarList(avatars)

		if len(serialized) != 2 || serialized[0].ID != "avatar-2" || serialized[1].ID != "avatar-1" {
			t.Fatalf("expected serialized avatars to preserve order, got %+v", serialized)
		}
	})

	t.Run("returns empty slice for empty input", func(t *testing.T) {
		serialized := avatardto.SerializeAvatarList(nil)

		if serialized == nil || len(serialized) != 0 {
			t.Fatalf("expected non-nil empty slice, got %+v", serialized)
		}
	})
}

func TestSerializeAvatarOptionsIncludesAvatarIDAndOptionSelection(t *testing.T) {
	t.Parallel()

	t.Run("maps avatar options", func(t *testing.T) {
		avatar := avatardomain.NewAvatarWithOptions(
			"avatar-v7",
			"user-v7",
			"Studio Hero",
			[]avatardomain.AvatarOption{
				{ID: "option-1", Href: "https://cdn.brandtoon.local/options/1.png", Selected: false},
				{ID: "option-2", Href: "https://cdn.brandtoon.local/options/2.png", Selected: true},
			},
		)

		serialized := avatardto.SerializeAvatarOptions(avatar)

		if serialized.AvatarID != "avatar-v7" || len(serialized.AvatarOptions) != 2 {
			t.Fatalf("expected avatar id and two options, got %+v", serialized)
		}

		if !serialized.AvatarOptions[1].Selected {
			t.Fatalf("expected selected option to remain selected, got %+v", serialized.AvatarOptions[1])
		}
	})

	t.Run("returns empty slice when avatar has no options", func(t *testing.T) {
		serialized := avatardto.SerializeAvatarOptions(avatardomain.NewAvatar("avatar-empty", "user-v7", "Empty"))

		if serialized.AvatarOptions == nil || len(serialized.AvatarOptions) != 0 {
			t.Fatalf("expected non-nil empty options slice, got %+v", serialized)
		}
	})
}

func TestSerializeAvatarDetailsIncludesNestedAvatarOptions(t *testing.T) {
	t.Parallel()

	avatar := avatardomain.NewAvatarWithOptions(
		"avatar-v7",
		"user-v7",
		"Studio Hero",
		[]avatardomain.AvatarOption{
			{ID: "option-1", Href: "https://cdn.brandtoon.local/options/1.png", Selected: true},
		},
	)

	serialized := avatardto.SerializeAvatarDetails(avatar)

	if serialized.ID != "avatar-v7" || serialized.Name != "Studio Hero" {
		t.Fatalf("expected serialized avatar details identity, got %+v", serialized)
	}

	if len(serialized.AvatarOptions) != 1 || serialized.AvatarOptions[0].ID != "option-1" {
		t.Fatalf("expected nested avatar options to be preserved, got %+v", serialized.AvatarOptions)
	}
}
