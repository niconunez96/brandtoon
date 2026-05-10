package avataroptiondomain

import (
	"errors"
	"strings"
)

type Status string

const (
	StatusPending Status = "PENDING"
	StatusDone    Status = "DONE"
	StatusFailed  Status = "FAILED"
)

var ErrAvatarOptionHrefRequired = errors.New("avatar option href required")
var ErrAvatarOptionOnlyDoneCanBeSelected = errors.New("only done avatar options can be selected")
var ErrAvatarOptionTerminalState = errors.New("avatar option already reached a terminal state")
var ErrAvatarOptionNotFound = errors.New("avatar option not found")

type AvatarOption struct {
	ID       string
	AvatarID string
	Status   Status
	Selected bool
	Href     *string
}

func NewPendingAvatarOption(id string, avatarID string) (AvatarOption, error) {
	option := AvatarOption{ID: id, AvatarID: avatarID, Status: StatusPending, Selected: false}
	return option, nil
}

func (o AvatarOption) MarkDone(href string) (AvatarOption, error) {
	if o.Status == StatusDone || o.Status == StatusFailed {
		return o, ErrAvatarOptionTerminalState
	}

	if strings.TrimSpace(href) == "" {
		return o, ErrAvatarOptionHrefRequired
	}

	nextHref := href
	o.Status = StatusDone
	o.Href = &nextHref
	return o, nil
}

func (o AvatarOption) MarkFailed() (AvatarOption, error) {
	if o.Status == StatusDone || o.Status == StatusFailed {
		return o, ErrAvatarOptionTerminalState
	}

	o.Status = StatusFailed
	o.Selected = false
	o.Href = nil
	return o, nil
}

func (o AvatarOption) Select() (AvatarOption, error) {
	if o.Status != StatusDone {
		return o, ErrAvatarOptionOnlyDoneCanBeSelected
	}

	o.Selected = true
	return o, nil
}
