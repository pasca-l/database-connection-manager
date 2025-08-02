package cli

import "errors"

var (
	ErrConfigAlreadyExists       = errors.New("configuration already exists")
	ErrUnsupportedConnectionType = errors.New("unsupported connection type")
	ErrDatabaseRequired          = errors.New("database name is required")
	ErrUsernameRequired          = errors.New("username is required")
)
