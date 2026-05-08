package avatardomain

import "errors"

var ErrInvalidName = errors.New("invalid avatar name")
var ErrAvatarOptionNotFound = errors.New("avatar option not found")

type AvatarOption struct {
	ID       string
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

func (a Avatar) WithNormalizedOptions() Avatar {
	a.AvatarOptions = normalizeAvatarOptions(a.AvatarOptions)
	return a
}

func (a Avatar) SelectOption(optionID string) (Avatar, error) {
	selectedIndex := -1
	for index, option := range a.AvatarOptions {
		if option.ID == optionID {
			selectedIndex = index
			break
		}
	}

	if selectedIndex == -1 {
		return a, ErrAvatarOptionNotFound
	}

	normalized := cloneAvatarOptions(a.AvatarOptions)
	for index := range normalized {
		normalized[index].Selected = index == selectedIndex
	}

	a.AvatarOptions = normalized
	return a, nil
}

func (a Avatar) DeleteOptions(optionIDs []string) (Avatar, error) {
	if len(optionIDs) == 0 {
		return a, ErrAvatarOptionNotFound
	}

	idsToDelete := make(map[string]struct{}, len(optionIDs))
	for _, optionID := range optionIDs {
		idsToDelete[optionID] = struct{}{}
	}

	remaining := make([]AvatarOption, 0, len(a.AvatarOptions))
	deletedCount := 0
	for _, option := range a.AvatarOptions {
		if _, shouldDelete := idsToDelete[option.ID]; shouldDelete {
			deletedCount++
			continue
		}

		remaining = append(remaining, option)
	}

	if deletedCount == 0 {
		return a, ErrAvatarOptionNotFound
	}

	a.AvatarOptions = normalizeAvatarOptions(remaining)
	return a, nil
}

func normalizeAvatarOptions(options []AvatarOption) []AvatarOption {
	if len(options) == 0 {
		return []AvatarOption{}
	}

	normalized := cloneAvatarOptions(options)
	selectedIndex := -1
	for index, option := range normalized {
		if !option.Selected {
			continue
		}

		selectedIndex = index
		break
	}

	if selectedIndex == -1 {
		selectedIndex = 0
	}

	for index := range normalized {
		normalized[index].Selected = index == selectedIndex
	}

	return normalized
}
