package avataroptionusecases

import "errors"

var ErrAvatarNotFound = errors.New("avatar not found")
var ErrAvatarConfigNotFound = errors.New("avatar config not found")
var ErrAvatarOptionNotFound = errors.New("avatar option not found")
var ErrAvatarGenerationUnavailable = errors.New("avatar generation unavailable")
