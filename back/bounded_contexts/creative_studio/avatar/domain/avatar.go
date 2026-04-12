package avatardomain

type Avatar struct {
	ID     string
	Name   string
	UserID string
}

func NewAvatar(id string, userID string, name string) Avatar {
	return Avatar{
		ID:     id,
		Name:   name,
		UserID: userID,
	}
}
