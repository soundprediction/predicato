package utils

import (
	"fmt"
	"log/slog"
	"runtime/debug"
)

// PanicError wraps a panic value as an error
type PanicError struct {
	Value      interface{}
	StackTrace string
}

func (e *PanicError) Error() string {
	return fmt.Sprintf("panic: %v", e.Value)
}

// RecoverAsError recovers from a panic and converts it to an error.
// It should be called with defer at the beginning of a function.
// The errPtr should be a pointer to the error return value.
//
// Example:
//
//	func doWork() (err error) {
//	    defer RecoverAsError(&err)
//	    // ... code that might panic
//	}
func RecoverAsError(errPtr *error) {
	if r := recover(); r != nil {
		stack := string(debug.Stack())
		*errPtr = &PanicError{
			Value:      r,
			StackTrace: stack,
		}
		slog.Error("Recovered from panic", "panic", r, "stack", stack)
	}
}

// RecoverWithCallback recovers from a panic and calls the callback with the error.
// Useful when you can't use the error return pattern.
//
// Example:
//
//	func doWork() {
//	    defer RecoverWithCallback(func(err error) {
//	        log.Printf("Work failed: %v", err)
//	    })
//	    // ... code that might panic
//	}
func RecoverWithCallback(callback func(error)) {
	if r := recover(); r != nil {
		stack := string(debug.Stack())
		err := &PanicError{
			Value:      r,
			StackTrace: stack,
		}
		slog.Error("Recovered from panic", "panic", r, "stack", stack)
		if callback != nil {
			callback(err)
		}
	}
}

// SafeGo runs a function in a goroutine with panic recovery.
// Any panic is logged and passed to the optional error handler.
//
// Example:
//
//	SafeGo(func() {
//	    // code that might panic
//	}, func(err error) {
//	    log.Printf("Goroutine failed: %v", err)
//	})
func SafeGo(fn func(), onError func(error)) {
	go func() {
		defer RecoverWithCallback(onError)
		fn()
	}()
}

// SafeGoWithContext runs a function in a goroutine with panic recovery.
// Returns a channel that receives any error (including panic errors).
// The channel is closed when the function completes.
//
// Example:
//
//	errCh := SafeGoWithContext(func() error {
//	    // code that might panic or return error
//	    return nil
//	})
//	if err := <-errCh; err != nil {
//	    log.Printf("Goroutine failed: %v", err)
//	}
func SafeGoWithResult(fn func() error) <-chan error {
	errCh := make(chan error, 1)
	go func() {
		defer close(errCh)
		defer func() {
			if r := recover(); r != nil {
				stack := string(debug.Stack())
				err := &PanicError{
					Value:      r,
					StackTrace: stack,
				}
				slog.Error("Recovered from panic in goroutine", "panic", r, "stack", stack)
				errCh <- err
			}
		}()
		if err := fn(); err != nil {
			errCh <- err
		}
	}()
	return errCh
}

// MustRecover is a helper for tests that ensures a panic is recovered
// and converted to an error for assertion.
func MustRecover(fn func()) (recovered interface{}) {
	defer func() {
		recovered = recover()
	}()
	fn()
	return nil
}
