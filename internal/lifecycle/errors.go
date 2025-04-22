package lifecycle

import "errors"

var (
	// ErrAlreadyRunning indicates that the component is already in a running state
	ErrAlreadyRunning = errors.New("component is already running")

	// ErrNotRunning indicates that the component is not in a running state
	ErrNotRunning = errors.New("component is not running")

	// ErrInvalidState indicates that the component is in an invalid state for the requested operation
	ErrInvalidState = errors.New("component is in an invalid state")

	// ErrInitializationFailed indicates that the component failed to initialize
	ErrInitializationFailed = errors.New("component initialization failed")

	// ErrShutdownFailed indicates that the component failed to shutdown properly
	ErrShutdownFailed = errors.New("component shutdown failed")
)
