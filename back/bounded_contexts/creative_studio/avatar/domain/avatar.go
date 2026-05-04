package avatardomain

import "errors"

var ErrInvalidName = errors.New("invalid avatar name")

type AvatarOption struct {
	Href     string
	Selected bool
}

type Avatar struct {
	ID            string
	Name          string
	UserID        string
	AvatarOptions []AvatarOption
}

func NewAvatar(id string, userID string, name string) Avatar {
	return NewAvatarWithoutOptions(id, userID, name)
}

func NewAvatarWithoutOptions(id string, userID string, name string) Avatar {
	return NewAvatarWithOptions(id, userID, name, nil)
}

func NewAvatarWithOptions(
	id string,
	userID string,
	name string,
	avatarOptions []AvatarOption,
) Avatar {
	return Avatar{
		ID:            id,
		Name:          name,
		UserID:        userID,
		AvatarOptions: cloneAvatarOptions(avatarOptions),
	}
}

func cloneAvatarOptions(options []AvatarOption) []AvatarOption {
	if len(options) == 0 {
		return []AvatarOption{}
	}

	cloned := make([]AvatarOption, len(options))
	copy(cloned, options)
	return cloned
}

func (a Avatar) AddOptions(options []AvatarOption) Avatar {
	combined := make([]AvatarOption, 0, len(a.AvatarOptions)+len(options))
	combined = append(combined, a.AvatarOptions...)
	combined = append(combined, options...)
	a.AvatarOptions = cloneAvatarOptions(combined)
	return a
}
