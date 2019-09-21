package relay

import (
	"fmt"
	"os"
	"path"
)

type Handler func(error)

type Relay struct {
	err error
	h   Handler
}

func New(hs ...Handler) *Relay {
	h := Handle
	if hs != nil && len(hs) > 0 {
		h = hs[0]
	}

	return &Relay{
		h: h,
	}
}

func (r *Relay) E(err error) {
	if err == nil {
		return
	}

	r.err = err

	panic(r)
}

func (r *Relay) Filter(v interface{}) {
	if v == nil {
		return
	}

	if rx, ok := v.(*Relay); !ok || r != rx {
		panic(v)
	}

	r.h(r.err)
}

func Handle(err error) {
	cmd := path.Base(os.Args[0])
	fmt.Fprintf(os.Stderr, "%s: %v\n", cmd, err)

	code := 1

	if coder, ok := err.(interface{ Code() int }); ok {
		code = coder.Code()
	}

	os.Exit(code)
}

type codedError struct {
	error
	code int
}

func (ce *codedError) Code() int {
	return ce.code
}

func CodedEFunc(r *Relay) func(error, int) {
	return func(err error, code int) {
		if err == nil {
			return
		}

		r.E(&codedError{err, code})
	}
}
