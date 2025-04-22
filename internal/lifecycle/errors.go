package lifecycle

import "errors"

var (
	// ErrNotRunning indicates that the component is not running when it should be
	ErrNotRunning = errors.New("component is not running")

	// ErrAlreadyRunning indicates that the component is already running
	ErrAlreadyRunning = errors.New("component is already running")

	// ErrInvalidState indicates that the component is in an invalid state
	ErrInvalidState = errors.New("component is in an invalid state")

	// ErrNotInitialized indicates that the component has not been initialized
	ErrNotInitialized = errors.New("component is not initialized")

	// ErrInitializationFailed indicates that the component failed to initialize
	ErrInitializationFailed = errors.New("component initialization failed")

	// ErrShutdownFailed indicates that the component failed to shutdown properly
	ErrShutdownFailed = errors.New("component shutdown failed")
)
