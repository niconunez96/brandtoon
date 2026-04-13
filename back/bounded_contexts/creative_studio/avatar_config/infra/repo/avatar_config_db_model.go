package avatarconfigrepo

import avatarconfigdomain "brandtoonapi/bounded_contexts/creative_studio/avatar_config/domain"

type avatarConfigDBModel struct {
	AvatarID      string `db:"avatar_id"`
	ArtisticStyle string `db:"artistic_style"`
	Prompt        string `db:"prompt"`
}

func newAvatarConfigDBModel(
	avatarConfig avatarconfigdomain.AvatarConfig,
) *avatarConfigDBModel {
	return &avatarConfigDBModel{
		AvatarID:      avatarConfig.AvatarID,
		ArtisticStyle: string(avatarConfig.ArtisticStyle),
		Prompt:        avatarConfig.Prompt,
	}
}

func (m *avatarConfigDBModel) ToDomain() (avatarconfigdomain.AvatarConfig, error) {
	artisticStyle, err := avatarconfigdomain.ParseArtisticStyle(m.ArtisticStyle)
	if err != nil {
		return avatarconfigdomain.AvatarConfig{}, err
	}

	return avatarconfigdomain.NewAvatarConfig(m.AvatarID, m.Prompt, artisticStyle), nil
}
