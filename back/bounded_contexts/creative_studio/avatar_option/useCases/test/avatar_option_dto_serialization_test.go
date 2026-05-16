package test

import (
	"testing"

	avataroptiondomain "brandtoonapi/bounded_contexts/creative_studio/avatar_option/domain"
	avataroptiondto "brandtoonapi/bounded_contexts/creative_studio/avatar_option/useCases/dto"
)

func TestSerializeMapsSingleAvatarOption(t *testing.T) {
	t.Parallel()

	href := "https://cdn.brandtoon.local/avatars/avatar-v7/options/1.png"
	option, err := avataroptiondomain.NewPendingAvatarOption("option-1", "avatar-v7")
	if err != nil {
		t.Fatalf("expected pending option, got %v", err)
	}
	option, err = option.MarkDone(href)
	if err != nil {
		t.Fatalf("expected done option, got %v", err)
	}

	serialized := avataroptiondto.Serialize(option)

	if serialized.ID != "option-1" || serialized.Status != string(avataroptiondomain.StatusDone) ||
		serialized.Href == nil ||
		*serialized.Href != href {
		t.Fatalf("expected serialized option fields, got %+v", serialized)
	}
}

func TestSerializeListWrapsAvatarOptionsWithAvatarID(t *testing.T) {
	t.Parallel()

	href := "https://cdn.brandtoon.local/avatars/avatar-v7/options/1.png"
	option, err := avataroptiondomain.NewPendingAvatarOption("option-1", "avatar-v7")
	if err != nil {
		t.Fatalf("expected pending option, got %v", err)
	}
	option, err = option.MarkDone(href)
	if err != nil {
		t.Fatalf("expected done option, got %v", err)
	}

	serialized := avataroptiondto.SerializeList("avatar-v7", []avataroptiondomain.AvatarOption{option})

	if serialized.AvatarID != "avatar-v7" || len(serialized.AvatarOptions) != 1 ||
		serialized.AvatarOptions[0].ID != "option-1" {
		t.Fatalf("expected wrapped serialized options, got %+v", serialized)
	}
}
