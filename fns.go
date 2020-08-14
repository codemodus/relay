package relay

import (
	"fmt"
	"os"
	"path"
)

type CheckFunc func(err error)

type CodedCheckFunc func(code int, err error)

type TripFunc func(format string, args ...interface{})

type CodedTripFunc func(code int, format string, args ...interface{})

// TripFn wraps the provided check function so that it can be called with a
// formatted error. This enables the immediate tripping of the related relay.
func TripFn(ck CheckFunc) TripFunc {
	return func(format string, args ...interface{}) {
		ck(fmt.Errorf(format, args...))
	}
}

// CodedTripFn wraps the provided codedCheck function so that it can be called
// with an exit code and formatted error. This enables the immediate tripping
// of the related relay.
func CodedTripFn(ck CodedCheckFunc) CodedTripFunc {
	return func(code int, format string, args ...interface{}) {
		ck(code, fmt.Errorf(format, args...))
	}
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

func Fns(handler ...func(error)) (CheckFunc, TripFunc) {
	r := New(handler...)
	c := r.Check
	t := TripFn(c)

	return c, t
}

func CodedFns(handler ...func(error)) (CodedCheckFunc, CodedTripFunc) {
	r := New(handler...)
	c := r.CodedCheck
	t := CodedTripFn(c)

	return c, t
}
