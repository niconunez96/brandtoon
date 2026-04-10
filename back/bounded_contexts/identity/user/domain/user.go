package domain

type User struct {
	AvatarURL string
	Email     string
	ID        string
	Name      string
}

func NewUser(id string, email string, name string, avatarURL string) User {
	return User{
		AvatarURL: avatarURL,
		Email:     email,
		ID:        id,
		Name:      name,
	}
}

func (u User) UpdateProfile(email string, name string, avatarURL string) User {
	u.Email = email
	u.Name = name
	u.AvatarURL = avatarURL
	return u
}
