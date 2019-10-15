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
		fmt.Println("extra message")
	}

	r := relay.New(h)
	defer func() { r.Filter(recover()) }()

	err := fail()
	r.Check(err)

	fmt.Println("should not print")

	// Output:
	// always fails
	// extra message
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

func ExampleTripFn() {
	check, filter := relay.New().Fns()
	trip := relay.TripFn(check)
	defer func() { filter(recover()) }()

	n := three()
	if n != 2 {
		trip("must receive %v: %v is invalid", 2, n)
	}

	fmt.Println("should not print")

	// prints "{cmd_name}: {trip_msg}" to stderr
	// calls os.Exit with code set as 1
}

func ExampleCodedTripFn() {
	check, filter := relay.New().CodedFns()
	trip := relay.CodedTripFn(check)
	defer func() { filter(recover()) }()

	n := three()
	if n != 2 {
		trip(4, "must receive %v: %v is invalid", 2, n)
	}

	fmt.Println("should not print")

	// prints "{cmd_name}: {trip_msg}" to stderr
	// calls os.Exit with code set as first arg to "trip"
}

func three() int {
	return 3
}

func fail() error {
	return errors.New("always fails")
}
