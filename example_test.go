package relay_test

import (
	"errors"
	"fmt"

	"github.com/codemodus/relay"
)

func Example() {
	r := relay.New()
	defer relay.Handle()

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
	defer relay.Handle()

	err := fail()
	r.Check(err)

	fmt.Println("should not print")

	// Output:
	// always fails
	// extra message
}

// Store check method for convenience.
func Example_easedUsage() {
	ck := relay.New().Check
	defer relay.Handle()

	err := fail()
	ck(err)

	// prints "{cmd_name}: {err_msg}" to stderr
	// calls os.Exit with code set as 1
}

func Example_exitCode() {
	ck := relay.New().CodedCheck
	defer relay.Handle()

	err := fail()
	ck(3, err)

	// prints "{cmd_name}: {err_msg}" to stderr
	// calls os.Exit with code set as first arg to r.CodedCheck ("ck")
}

func ExampleTripFn() {
	ck := relay.New().Check
	trip := relay.TripFn(ck)
	defer relay.Handle()

	n := three()
	if n != 2 {
		trip("must receive %v: %v is invalid", 2, n)
	}

	fmt.Println("should not print")

	// prints "{cmd_name}: {trip_msg}" to stderr
	// calls os.Exit with code set as 1
}

func ExampleCodedTripFn() {
	ck := relay.New().CodedCheck
	trip := relay.CodedTripFn(ck)
	defer relay.Handle()

	n := three()
	if n != 2 {
		trip(4, "must receive %v: %v is invalid", 2, n)
	}

	fmt.Println("should not print")

	// prints "{cmd_name}: {trip_msg}" to stderr
	// calls os.Exit with code set as first arg to "trip"
}

func ExampleFns() {
	ck, trip := relay.Fns()
	defer relay.Handle()

	err := mightFail()
	ck(err)

	n := three()
	if n != 2 {
		trip("must receive %v: %v is invalid", 2, n)
	}

	fmt.Println("should not print")

	// prints "{cmd_name}: {trip_msg}" to stderr
	// calls os.Exit with code set as 1
}

func ExampleCodedFns() {
	ck, trip := relay.CodedFns()
	defer relay.Handle()

	err := mightFail()
	ck(3, err)

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

func mightFail() error {
	return nil
}
