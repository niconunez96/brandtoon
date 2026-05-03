package avatarconfigrepo

import avatarconfigdomain "brandtoonapi/bounded_contexts/creative_studio/avatar_config/domain"

type avatarConfigDBModel struct {
	AvatarID      string `db:"avatar_id"`
	ArtisticStyle string `db:"artistic_style"`
	Personality   string `db:"personality"`
	Prompt        string `db:"prompt"`
}

func newAvatarConfigDBModel(
	avatarConfig avatarconfigdomain.AvatarConfig,
) *avatarConfigDBModel {
	return &avatarConfigDBModel{
		AvatarID:      avatarConfig.AvatarID,
		ArtisticStyle: string(avatarConfig.ArtisticStyle),
		Personality:   string(avatarConfig.Personality),
		Prompt:        avatarConfig.Prompt,
	}
}

func (m *avatarConfigDBModel) ToDomain() (avatarconfigdomain.AvatarConfig, error) {
	artisticStyle, err := avatarconfigdomain.ParseArtisticStyle(m.ArtisticStyle)
	if err != nil {
		return avatarconfigdomain.AvatarConfig{}, err
	}

	personality, err := avatarconfigdomain.ParsePersonality(m.Personality)
	if err != nil {
		return avatarconfigdomain.AvatarConfig{}, err
	}

	return avatarconfigdomain.NewAvatarConfig(
		m.AvatarID,
		m.Prompt,
		artisticStyle,
		personality,
	), nil
}
