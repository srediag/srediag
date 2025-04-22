package lifecycle

import (
	"context"
	"testing"
	"time"
)

func TestBaseManager_InitialState(t *testing.T) {
	bm := NewBaseManager()
	if bm.IsRunning() {
		t.Error("novo manager não deve estar rodando")
	}
	if !bm.IsHealthy() {
		t.Error("novo manager deve estar saudável")
	}
}

func TestBaseManager_Start(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name      string
		setup     func(*BaseManager)
		wantError bool
	}{
		{
			name:      "start quando parado",
			setup:     func(bm *BaseManager) {},
			wantError: false,
		},
		{
			name: "start quando já rodando",
			setup: func(bm *BaseManager) {
				bm.SetRunning(true)
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bm := NewBaseManager()
			tt.setup(bm)

			err := bm.Start(ctx)
			if (err != nil) != tt.wantError {
				t.Errorf("Start() error = %v, wantError %v", err, tt.wantError)
			}

			if !tt.wantError && !bm.IsRunning() {
				t.Error("Start() deve definir o estado como rodando")
			}
		})
	}
}

func TestBaseManager_Stop(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name      string
		setup     func(*BaseManager)
		wantError bool
	}{
		{
			name: "stop quando rodando",
			setup: func(bm *BaseManager) {
				bm.SetRunning(true)
			},
			wantError: false,
		},
		{
			name:      "stop quando já parado",
			setup:     func(bm *BaseManager) {},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bm := NewBaseManager()
			tt.setup(bm)

			err := bm.Stop(ctx)
			if (err != nil) != tt.wantError {
				t.Errorf("Stop() error = %v, wantError %v", err, tt.wantError)
			}

			if !tt.wantError && bm.IsRunning() {
				t.Error("Stop() deve definir o estado como não rodando")
			}
		})
	}
}

func TestBaseManager_Health(t *testing.T) {
	bm := NewBaseManager()

	// Teste estado inicial
	if !bm.IsHealthy() {
		t.Error("estado inicial de saúde deve ser true")
	}

	// Teste atualização de estado
	bm.UpdateHealth(false)
	if bm.IsHealthy() {
		t.Error("saúde deve ser false após UpdateHealth(false)")
	}

	bm.UpdateHealth(true)
	if !bm.IsHealthy() {
		t.Error("saúde deve ser true após UpdateHealth(true)")
	}
}

func TestBaseManager_ConcurrentAccess(t *testing.T) {
	bm := NewBaseManager()
	done := make(chan bool)
	iterations := 1000

	// Teste acesso concorrente ao estado running
	go func() {
		for i := 0; i < iterations; i++ {
			bm.SetRunning(true)
			bm.SetRunning(false)
		}
		done <- true
	}()

	go func() {
		for i := 0; i < iterations; i++ {
			bm.IsRunning()
		}
		done <- true
	}()

	// Teste acesso concorrente ao estado de saúde
	go func() {
		for i := 0; i < iterations; i++ {
			bm.UpdateHealth(true)
			bm.UpdateHealth(false)
		}
		done <- true
	}()

	go func() {
		for i := 0; i < iterations; i++ {
			bm.IsHealthy()
		}
		done <- true
	}()

	// Adiciona timeout para evitar deadlock
	timeout := time.After(5 * time.Second)
	for i := 0; i < 4; i++ {
		select {
		case <-done:
			continue
		case <-timeout:
			t.Fatal("timeout esperando goroutines terminarem")
		}
	}
}
