package avatarrepo

import (
	avatardomain "brandtoonapi/bounded_contexts/creative_studio/avatar/domain"
	"testing"
)

func TestEncodeAvatarOptionsJSONReturnsStableEmptyArray(t *testing.T) {
	t.Parallel()

	encoded, err := encodeAvatarOptionsJSON(nil)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if string(encoded) != "[]" {
		t.Fatalf("expected [] payload, got %s", string(encoded))
	}
}

func TestDecodeAvatarOptionsJSONReturnsStructuredOptions(t *testing.T) {
	t.Parallel()

	decoded, err := decodeAvatarOptionsJSON([]byte(`[{"href":"https://cdn.brandtoon.local/a.png","selected":false}]`))
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if len(decoded) != 1 {
		t.Fatalf("expected 1 option, got %d", len(decoded))
	}

	if decoded[0] != (avatardomain.AvatarOption{Href: "https://cdn.brandtoon.local/a.png", Selected: false}) {
		t.Fatalf("unexpected option decoded: %+v", decoded[0])
	}
}

func TestAvatarConfigUpsertResponsePreservesAvatarOwnedOptionsBoundary(t *testing.T) {
	t.Parallel()

	model := &avatarDBModel{
		ID:                "avatar-v7",
		Name:              "Studio Hero",
		UserID:            "user-v7",
		AvatarOptionsJSON: []byte(`[{"href":"https://cdn.brandtoon.local/a.png","selected":false}]`),
	}

	avatar := model.ToDomain()
	if len(avatar.AvatarOptions) != 1 {
		t.Fatalf("expected avatar aggregate to own options, got %d", len(avatar.AvatarOptions))
	}
}
