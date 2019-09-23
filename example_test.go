package relay_test

import (
	"errors"
	"fmt"

	"github.com/codemodus/relay"
)

func Example() {
	r := relay.New()
	defer func() { r.Filter(recover()) }()

	err := fail()
	r.Check(err)
	// prints "{cmd_name}: {err_msg}" to stderr
	// calls os.Exit with code set as 1
}

// Override default error handler.
func Example_customHandler() {
	h := func(err error) {
		fmt.Println(err)
	}

	r := relay.New(h)
	defer func() { r.Filter(recover()) }()

	defer fmt.Println("reached")

	err := fail()
	r.Check(err)

	fmt.Println("should not print")

	// Output:
	// reached
	// always fails
}

// Convenience methods for eased usage.
func Example_easedUsage() {
	ce, filter := relay.New().Fns()
	defer func() { filter(recover()) }()

	err := fail()
	ce(err)
	// prints "{cmd_name}: {err_msg}" to stderr
	// calls os.Exit with code set as 1
}

func ExampleRelay_CodedCheck() {
	r := relay.New()
	defer func() { r.Filter(recover()) }()

	err := fail()
	r.CodedCheck(3, err)
	// prints "{cmd_name}: {err_msg}" to stderr
	// calls os.Exit with code set as first arg to r.CodedCheck
}

func ExampleRelay_Fns() {
	ce, filter := relay.New().Fns()
	defer func() { filter(recover()) }()

	err := fail()
	ce(err)
	// prints "{cmd_name}: {err_msg}" to stderr
	// calls os.Exit with code set as 1
}

func ExampleRelay_CodedFns() {
	ce, filter := relay.New().CodedFns()
	defer func() { filter(recover()) }()

	err := fail()
	ce(3, err)
	// prints "{cmd_name}: {err_msg}" to stderr
	// calls os.Exit with code set as first arg to ce
}

func ExampleCoder() {
	r := relay.New()
	defer func() { r.Filter(recover()) }()

	err := fail()
	cerr := &relay.CodedError{err, 2} // satisfies the Coder interface
	r.Check(cerr)
	// prints "{cmd_name}: {err_msg}" to stderr
	// calls os.Exit with code set as cerr.Code() return value
}

func fail() error {
	return errors.New("always fails")
}
