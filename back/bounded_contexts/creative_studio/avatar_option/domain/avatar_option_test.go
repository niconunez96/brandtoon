package avataroptiondomain

import (
	"errors"
	"testing"
)

func TestNewPendingAvatarOptionStartsPendingWithoutHref(t *testing.T) {
	t.Parallel()

	option, err := NewPendingAvatarOption("option-v7", "avatar-v7")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if option.Status != StatusPending {
		t.Fatalf("expected pending status, got %s", option.Status)
	}

	if option.Href != nil {
		t.Fatalf("expected nil href for pending option, got %v", option.Href)
	}

	if option.Selected {
		t.Fatalf("expected pending option not selected")
	}
}

func TestAvatarOptionMarkDoneRequiresHref(t *testing.T) {
	t.Parallel()

	option, err := NewPendingAvatarOption("option-v7", "avatar-v7")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	updated, err := option.MarkDone("https://cdn.brandtoon.local/avatars/avatar-v7/options/1.png")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if updated.Status != StatusDone {
		t.Fatalf("expected done status, got %s", updated.Status)
	}

	if updated.Href == nil || *updated.Href == "" {
		t.Fatalf("expected href for done option, got %+v", updated)
	}

	_, err = option.MarkDone("")
	if !errors.Is(err, ErrAvatarOptionHrefRequired) {
		t.Fatalf("expected ErrAvatarOptionHrefRequired, got %v", err)
	}
}

func TestAvatarOptionSelectRequiresDoneStatus(t *testing.T) {
	t.Parallel()

	pending, err := NewPendingAvatarOption("option-v7", "avatar-v7")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	_, err = pending.Select()
	if !errors.Is(err, ErrAvatarOptionOnlyDoneCanBeSelected) {
		t.Fatalf("expected ErrAvatarOptionOnlyDoneCanBeSelected, got %v", err)
	}

	done, err := pending.MarkDone("https://cdn.brandtoon.local/avatars/avatar-v7/options/1.png")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	selected, err := done.Select()
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if !selected.Selected {
		t.Fatalf("expected selected option, got %+v", selected)
	}
}

func TestAvatarOptionTerminalStatesCannotTransitionAgain(t *testing.T) {
	t.Parallel()

	pending, err := NewPendingAvatarOption("option-v7", "avatar-v7")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	failed, err := pending.MarkFailed()
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	_, err = failed.MarkDone("https://cdn.brandtoon.local/avatars/avatar-v7/options/1.png")
	if !errors.Is(err, ErrAvatarOptionTerminalState) {
		t.Fatalf("expected ErrAvatarOptionTerminalState, got %v", err)
	}

	done, err := pending.MarkDone("https://cdn.brandtoon.local/avatars/avatar-v7/options/1.png")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	_, err = done.MarkFailed()
	if !errors.Is(err, ErrAvatarOptionTerminalState) {
		t.Fatalf("expected ErrAvatarOptionTerminalState, got %v", err)
	}
}
