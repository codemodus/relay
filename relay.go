// Package relay provides a simple mechanism for relaying control flow based
// upon whether a checked error is nil or not. This mechanism requires special
// setup within an application due to the behavior of the builtin function
// recover(). Please review the provided examples to ensure correct usage.
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
	h := Handle
	if handler != nil && len(handler) > 0 {
		h = handler[0]
	}

	return &Relay{
		h: h,
	}
}

// Check will do nothing if the error argument is nil. Otherwise, it kicks-off
// an event (i.e. the relay is "tripped") and should be handled by a deferred
// and wrapped call to r.Filter(recover()).
func (r *Relay) Check(err error) {
	if err == nil {
		return
	}

	r.err = err

	panic(r)
}

// CodedCheck will do nothing if the error argument is nil. Otherwise, it
// kicks-off an event (i.e. the relay is "tripped") and should be handled by a
// deferred and wrapped call to r.Filter(recover()). Any provided error will be
// wrapped in a CodedError instance in order to trigger special behavior in the
// default error handler.
func (r *Relay) CodedCheck(code int, err error) {
	if err == nil {
		return
	}

	r.Check(&CodedError{err, code})
}

// Filter should be wrapped and deferred before any usage of the receiver
// occurs. The argument should be a call to the recover() builtin. If no panic
// has been triggered, Filter will do nothing. Otherwise, it will handle the
// currently set err. If the argument value is not recognized, the value is
// passed into an additional panic() call.
func (r *Relay) Filter(v interface{}) {
	if v == nil {
		return
	}

	if rx, ok := v.(*Relay); !ok || r != rx {
		panic(v)
	}

	r.h(r.err)
}

// Coder describes any type that can return an error code.
type Coder interface {
	Code() int
}

// Handle is the default error handler. It will print "{cmd_name}: {err_msg}"
// to stderr and then call os.Exit. If the handled error happens to satisfy the
// Coder interface, that value will be used as the exit code. Otherwise, 1 will
// be used.
func Handle(err error) {
	if err == nil {
		return
	}

	cmd := path.Base(os.Args[0])
	fmt.Fprintf(os.Stderr, "%s: %v\n", cmd, err)

	code := 1
	if c, ok := err.(Coder); ok {
		code = c.Code()
	}

	os.Exit(code)
}

// CodedError is a simple implementaion of the Coder interface.
type CodedError struct {
	Err error
	C   int
}

// Error satisfies the error interface.
func (ce *CodedError) Error() string {
	return ce.Err.Error()
}

// Code satisfies the Coder interface.
func (ce *CodedError) Code() int {
	return ce.C
}

// Fns is a convenience method which returns the Check and Filter functions.
func (r *Relay) Fns() (check func(error), filter func(interface{})) {
	return r.Check, r.Filter
}

// CodedFns is a convenience method which returns the CodedCheck and Filter
// functions.
func (r *Relay) CodedFns() (codedCheck func(int, error), filter func(interface{})) {
	return r.CodedCheck, r.Filter
}

// TripFn wraps the provided check function so that it can be called with a
// formatted error. This enables the immediate triggering of the related relay.
func TripFn(check func(error)) func(string, ...interface{}) {
	return func(format string, args ...interface{}) {
		check(fmt.Errorf(format, args...))
	}
}

// CodedTripFn wraps the provided codedCheck function so that it can be called
// with an exit code and formatted error. This enables the immediate triggering
// of the related relay.
func CodedTripFn(codedCheck func(int, error)) func(int, string, ...interface{}) {
	return func(code int, format string, args ...interface{}) {
		codedCheck(code, fmt.Errorf(format, args...))
	}
}
