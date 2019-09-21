package relay_test

import (
	"errors"

	"github.com/codemodus/relay"
)

func Example() {
	r := relay.New()
	defer func() { r.Filter(recover()) }()
	re := r.E

	err := errors.New("always fails")
	re(err)
}

func ExampleCodedEFunc() {
	r := relay.New()
	defer func() { r.Filter(recover()) }()
	re := relay.CodedEFunc(r)

	err := errors.New("always fails")
	re(err, 42)
}
