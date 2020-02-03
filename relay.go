// Package relay provides a simple mechanism for relaying control flow based
// upon whether a checked error is nil or not.
package relay

import (
	"fmt"
	"os"
	"path"
)

// Relay tracks an error and a related error handler.
type Relay struct {
	err error
	h   func(error)
}

// New constructs a new *Relay.
func New(handler ...func(error)) *Relay {
	h := DefaultHandler()
	if len(handler) > 0 {
		h = handler[0]
	}

	return &Relay{
		h: h,
	}
}

// Check will do nothing if the error argument is nil. Otherwise, it kicks-off
// an event (i.e. the relay is "tripped") and should be handled by a deferred
// call to Handle().
func (r *Relay) Check(err error) {
	if err == nil {
		return
	}

	r.err = err

	panic(r)
}

// CodedCheck will do nothing if the error argument is nil. Otherwise, it
// kicks-off an event (i.e. the relay is "tripped") and should be handled by a
// deferred call to Handle(). Any provided error will be wrapped in a
// codedError instance in order to trigger special behavior in the default
// error handler.
func (r *Relay) CodedCheck(code int, err error) {
	if err == nil {
		return
	}

	r.Check(&CodedError{err, code})
}

// TripFn wraps the provided check function so that it can be called with a
// formatted error. This enables the immediate tripping of the related relay.
func TripFn(check func(error)) func(string, ...interface{}) {
	return func(format string, args ...interface{}) {
		check(fmt.Errorf(format, args...))
	}
}

// CodedTripFn wraps the provided codedCheck function so that it can be called
// with an exit code and formatted error. This enables the immediate tripping
// of the related relay.
func CodedTripFn(codedCheck func(int, error)) func(int, string, ...interface{}) {
	return func(code int, format string, args ...interface{}) {
		codedCheck(code, fmt.Errorf(format, args...))
	}
}

// ExitCoder describes any type that can return an error code.
type ExitCoder interface {
	ExitCode() int
}

// DefaultHandler returns an error handler that prints "{cmd_name}: {err_msg}"
// to stderr and then call os.Exit. If the handled error happens to satisfy the
// ExitCoder interface, that value will be used as the exit code. Otherwise, 1
// will be used.
func DefaultHandler() func(error) {
	return func(err error) {
		if err == nil {
			return
		}

		cmd := path.Base(os.Args[0])
		fmt.Fprintf(os.Stderr, "%s: %v\n", cmd, err)

		code := 1
		if ec, ok := err.(ExitCoder); ok {
			code = ec.ExitCode()
		}

		os.Exit(code)
	}
}

// Handle checks the recover() builtin and handles the error which tripped the
// relay, if any.
func Handle() {
	v := recover()
	if v == nil {
		return
	}

	r, ok := v.(*Relay)
	if !ok {
		panic(v)
	}

	r.h(r.err)
}

// CodedError is a simple implementaion of the ExitCoder interface.
type CodedError struct {
	Err error
	C   int
}

// Error satisfies the error interface.
func (ce *CodedError) Error() string {
	return ce.Err.Error()
}

// ExitCode satisfies the ExitCoder interface.
func (ce *CodedError) ExitCode() int {
	return ce.C
}
