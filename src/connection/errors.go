package connection

import "errors"

var (
	ErrUnsupportedConnectionType = errors.New("unsupported connection type")
	ErrConnectionAlreadyExists   = errors.New("connection already exists")
	ErrConnectionNotFound        = errors.New("connection not found")
	ErrBuildCommandFailed        = errors.New("failed to build command for connection")
)
