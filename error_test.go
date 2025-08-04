package oops

import (
	"errors"
	"io/fs"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorsIs(t *testing.T) {
	is := assert.New(t)

	err := Errorf("Error: %w", fs.ErrExist)
	is.True(errors.Is(err, fs.ErrExist))

	err = Wrap(fs.ErrExist)
	is.True(errors.Is(err, fs.ErrExist))

	err = Wrapf(fs.ErrExist, "Error: %w", assert.AnError)
	is.True(errors.Is(err, fs.ErrExist))

	err = Join(fs.ErrExist, assert.AnError)
	is.True(errors.Is(err, fs.ErrExist))
	err = Join(assert.AnError, fs.ErrExist)
	is.True(errors.Is(err, fs.ErrExist))

	err = Recover(func() {
		panic(fs.ErrExist)
	})
	is.True(errors.Is(err, fs.ErrExist))

	err = Recoverf(func() {
		panic(fs.ErrExist)
	}, "Error: %w", assert.AnError)
	is.True(errors.Is(err, fs.ErrExist))
}

func TestErrorsAs(t *testing.T) {
	is := assert.New(t)

	var anError error = &fs.PathError{Err: fs.ErrExist}
	var target *fs.PathError

	err := Errorf("error: %w", anError)
	is.True(errors.As(err, &target))

	err = Wrap(anError)
	is.True(errors.As(err, &target))

	err = Wrapf(anError, "Error: %w", assert.AnError)
	is.True(errors.As(err, &target))

	err = Join(anError, assert.AnError)
	is.True(errors.As(err, &target))
	err = Join(assert.AnError, anError)
	is.True(errors.As(err, &target))

	err = Recover(func() {
		panic(anError)
	})
	is.True(errors.As(err, &target))

	err = Recoverf(func() {
		panic(anError)
	}, "Error: %w", assert.AnError)
	is.True(errors.As(err, &target))
}

// TestErrorsIsWithOopsErrors tests the fix for comparing wrapped oops errors without panics
func TestErrorsIsWithOopsErrors(t *testing.T) {
	is := assert.New(t)

	// Test case 1: Compare wrapped oops error with original oops error
	// This was the main use case that was causing panics
	originalErr := New("user not found")
	wrappedErr := In("User/GetByEmail").Wrap(originalErr)
	
	// This should not panic and should return true
	is.True(errors.Is(wrappedErr, originalErr))
	
	// Test case 2: Compare multiple levels of wrapping
	level1Err := New("database error")
	level2Err := In("Repository").Wrap(level1Err)
	level3Err := In("Service").Wrap(level2Err)
	
	is.True(errors.Is(level3Err, level1Err))
	is.True(errors.Is(level3Err, level2Err))
	
	// Test case 3: Compare with different oops errors (should return false)
	differentErr := New("different error")
	is.False(errors.Is(wrappedErr, differentErr))
	
	// Test case 4: Compare wrapped oops error with standard error
	standardErr := errors.New("user not found")
	wrappedStandardErr := In("User/GetByEmail").Wrap(standardErr)
	is.True(errors.Is(wrappedStandardErr, standardErr))
	
	// Test case 5: Compare oops error with wrapped standard error
	oopsErr := New("user not found")
	wrappedStandardErr2 := In("User/GetByEmail").Wrap(standardErr)
	is.False(errors.Is(wrappedStandardErr2, oopsErr)) // Different error types
	
	// Test case 6: Test with nil errors
	is.False(errors.Is(wrappedErr, nil))
	is.False(errors.Is(nil, originalErr))
	
	// Test case 7: Test with same message but different instances
	// Note: New() creates different underlying error instances, so they shouldn't match
	err1 := New("same message")
	err2 := New("same message")
	is.False(errors.Is(err1, err2)) // Different instances, should not match
	
	// Test case 8: Test with different messages
	err3 := New("different message")
	is.False(errors.Is(err1, err3))
}
