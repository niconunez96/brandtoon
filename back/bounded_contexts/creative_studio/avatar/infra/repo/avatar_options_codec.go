package avatarrepo

import (
	avatardomain "brandtoonapi/bounded_contexts/creative_studio/avatar/domain"
	"encoding/json"
)

func encodeAvatarOptionsJSON(options []avatardomain.AvatarOption) ([]byte, error) {
	if len(options) == 0 {
		return []byte("[]"), nil
	}

	return json.Marshal(options)
}

func decodeAvatarOptionsJSON(payload []byte) ([]avatardomain.AvatarOption, error) {
	if len(payload) == 0 {
		return []avatardomain.AvatarOption{}, nil
	}

	options := make([]avatardomain.AvatarOption, 0)
	if err := json.Unmarshal(payload, &options); err != nil {
		return nil, err
	}

	return options, nil
}
