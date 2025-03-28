package commandcontext_test

import (
	"context"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/motain/of-catalog/internal/utils/commandcontext"
	"github.com/stretchr/testify/assert"
)

func TestHandleSignal(t *testing.T) {
	tests := []struct {
		name          string
		signalToSend  os.Signal
		expectedError bool
	}{
		{
			name:          "Handle SIGTERM signal",
			signalToSend:  syscall.SIGTERM,
			expectedError: false,
		},
		{
			name:          "Handle SIGINT signal",
			signalToSend:  os.Interrupt,
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			go commandcontext.HandleSignal(cancel)

			// Simulate sending the signal in a separate goroutine to avoid stopping the test
			go func() {
				time.Sleep(100 * time.Millisecond) // Give some time for the signal handler to start
				syscall.Kill(syscall.Getpid(), tt.signalToSend.(syscall.Signal))
			}()

			// Wait for the context to be canceled
			select {
			case <-ctx.Done():
				if ctx.Err() != context.Canceled {
					t.Errorf("unexpected context error: %v", ctx.Err())
				}
			case <-time.After(1 * time.Second):
				t.Fatal("context was not canceled in time")
			}
		})
	}
}
func TestInit(t *testing.T) {
	tests := []struct {
		name          string
		signalToSend  os.Signal
		expectedError bool
	}{
		{
			name:          "Init handles SIGTERM signal",
			signalToSend:  syscall.SIGTERM,
			expectedError: false,
		},
		{
			name:          "Init handles SIGINT signal",
			signalToSend:  os.Interrupt,
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := commandcontext.Init()

			// Simulate sending the signal in a separate goroutine to avoid stopping the test
			go func() {
				time.Sleep(100 * time.Millisecond) // Give some time for the signal handler to start
				syscall.Kill(syscall.Getpid(), tt.signalToSend.(syscall.Signal))
			}()

			// Wait for the context to be canceled
			select {
			case <-ctx.Done():
				assert.Equal(t, context.Canceled, ctx.Err(), "unexpected context error")
			case <-time.After(1 * time.Second):
				t.Fatal("context was not canceled in time")
			}
		})
	}
}
