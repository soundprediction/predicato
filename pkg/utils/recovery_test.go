package utils

import (
	"errors"
	"testing"
	"time"
)

func TestRecoverAsError(t *testing.T) {
	t.Run("recovers from panic", func(t *testing.T) {
		fn := func() (err error) {
			defer RecoverAsError(&err)
			panic("test panic")
		}

		err := fn()
		if err == nil {
			t.Fatal("expected error from panic recovery")
		}

		var panicErr *PanicError
		if !errors.As(err, &panicErr) {
			t.Fatalf("expected PanicError, got %T", err)
		}

		if panicErr.Value != "test panic" {
			t.Errorf("expected panic value 'test panic', got %v", panicErr.Value)
		}

		if panicErr.StackTrace == "" {
			t.Error("expected stack trace to be populated")
		}
	})

	t.Run("no error when no panic", func(t *testing.T) {
		fn := func() (err error) {
			defer RecoverAsError(&err)
			return nil
		}

		err := fn()
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("preserves original error", func(t *testing.T) {
		originalErr := errors.New("original error")
		fn := func() (err error) {
			defer RecoverAsError(&err)
			return originalErr
		}

		err := fn()
		if err != originalErr {
			t.Errorf("expected original error, got %v", err)
		}
	})
}

func TestRecoverWithCallback(t *testing.T) {
	t.Run("calls callback on panic", func(t *testing.T) {
		var capturedErr error
		fn := func() {
			defer RecoverWithCallback(func(err error) {
				capturedErr = err
			})
			panic("callback test")
		}

		fn()

		if capturedErr == nil {
			t.Fatal("expected callback to be called with error")
		}

		var panicErr *PanicError
		if !errors.As(capturedErr, &panicErr) {
			t.Fatalf("expected PanicError, got %T", capturedErr)
		}
	})

	t.Run("handles nil callback", func(t *testing.T) {
		fn := func() {
			defer RecoverWithCallback(nil)
			panic("nil callback test")
		}

		// Should not panic
		fn()
	})
}

func TestSafeGo(t *testing.T) {
	t.Run("executes function without panic", func(t *testing.T) {
		done := make(chan struct{})
		SafeGo(func() {
			close(done)
		}, nil)

		select {
		case <-done:
			// Success
		case <-time.After(time.Second):
			t.Fatal("function did not complete")
		}
	})

	t.Run("recovers from panic and calls error handler", func(t *testing.T) {
		errCh := make(chan error, 1)
		SafeGo(func() {
			panic("safe go panic")
		}, func(err error) {
			errCh <- err
		})

		select {
		case err := <-errCh:
			if err == nil {
				t.Fatal("expected error from panic")
			}
			var panicErr *PanicError
			if !errors.As(err, &panicErr) {
				t.Fatalf("expected PanicError, got %T", err)
			}
		case <-time.After(time.Second):
			t.Fatal("error handler was not called")
		}
	})
}

func TestSafeGoWithResult(t *testing.T) {
	t.Run("returns nil error on success", func(t *testing.T) {
		errCh := SafeGoWithResult(func() error {
			return nil
		})

		select {
		case err, ok := <-errCh:
			if ok && err != nil {
				t.Errorf("expected no error, got %v", err)
			}
		case <-time.After(time.Second):
			t.Fatal("channel was not closed")
		}
	})

	t.Run("returns error from function", func(t *testing.T) {
		expectedErr := errors.New("function error")
		errCh := SafeGoWithResult(func() error {
			return expectedErr
		})

		select {
		case err := <-errCh:
			if err != expectedErr {
				t.Errorf("expected %v, got %v", expectedErr, err)
			}
		case <-time.After(time.Second):
			t.Fatal("error was not returned")
		}
	})

	t.Run("returns panic error", func(t *testing.T) {
		errCh := SafeGoWithResult(func() error {
			panic("result panic")
		})

		select {
		case err := <-errCh:
			if err == nil {
				t.Fatal("expected error from panic")
			}
			var panicErr *PanicError
			if !errors.As(err, &panicErr) {
				t.Fatalf("expected PanicError, got %T", err)
			}
		case <-time.After(time.Second):
			t.Fatal("error was not returned")
		}
	})
}

func TestMustRecover(t *testing.T) {
	t.Run("captures panic value", func(t *testing.T) {
		recovered := MustRecover(func() {
			panic("must recover test")
		})

		if recovered != "must recover test" {
			t.Errorf("expected 'must recover test', got %v", recovered)
		}
	})

	t.Run("returns nil when no panic", func(t *testing.T) {
		recovered := MustRecover(func() {
			// No panic
		})

		if recovered != nil {
			t.Errorf("expected nil, got %v", recovered)
		}
	})
}

func TestPanicErrorString(t *testing.T) {
	err := &PanicError{Value: "test value"}
	expected := "panic: test value"
	if err.Error() != expected {
		t.Errorf("expected %q, got %q", expected, err.Error())
	}
}

func TestConcurrentPanicRecovery(t *testing.T) {
	// Test that multiple concurrent goroutines with panics are all recovered
	// using SafeGoWithResult which properly handles the panic/error flow
	const numGoroutines = 10
	errChannels := make([]<-chan error, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		i := i
		errChannels[i] = SafeGoWithResult(func() error {
			if i%2 == 0 {
				panic("even panic")
			}
			return nil
		})
	}

	// Wait for all goroutines and count errors
	errorCount := 0
	for _, ch := range errChannels {
		select {
		case err := <-ch:
			if err != nil {
				errorCount++
			}
		case <-time.After(5 * time.Second):
			t.Fatal("timed out waiting for goroutine")
		}
	}

	// Half of the goroutines should have panicked
	if errorCount != numGoroutines/2 {
		t.Errorf("expected %d errors, got %d", numGoroutines/2, errorCount)
	}
}
