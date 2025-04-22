package core

import "errors"

var (
	// ErrAlreadyRunning indica que o componente já está em execução
	ErrAlreadyRunning = errors.New("component is already running")

	// ErrNotRunning indica que o componente não está em execução
	ErrNotRunning = errors.New("component is not running")

	// ErrInvalidConfig indica uma configuração inválida
	ErrInvalidConfig = errors.New("invalid configuration")

	// ErrInitializationFailed indica falha na inicialização
	ErrInitializationFailed = errors.New("initialization failed")

	// ErrShutdownFailed indica falha no desligamento
	ErrShutdownFailed = errors.New("shutdown failed")
)
