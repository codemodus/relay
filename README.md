# relay

    go get github.com/codemodus/relay

Package relay provides a simple mechanism for relaying control flow based
upon whether a checked error is nil or not. This mechanism requires special
setup within an application due to the behavior of the builtin function
recover(). Please review the provided examples to ensure correct usage.

## Usage

```go
func CodedTripFn(check func(int, error)) func(int, string, ...interface{})
func Handle(err error)
func TripFn(check func(error)) func(string, ...interface{})
type CodedError
    func (ce *CodedError) Code() int
    func (ce *CodedError) Error() string
type Coder
type Relay
    func New(handler ...func(error)) *Relay
    func (r *Relay) Check(err error)
    func (r *Relay) CodedCheck(code int, err error)
    func (r *Relay) CodedFns() (codedCheck func(int, error), filter func(interface{}))
    func (r *Relay) Filter(v interface{})
    func (r *Relay) Fns() (check func(error), filter func(interface{}))
```

### Setup

```go
import (
    "github.com/codemodus/relay"
)

func main() {
    r := relay.New()
    defer func() { r.Filter(recover()) }()

    err := fail()
    r.Check(err)
    // prints "{cmd_name}: {err_msg}" to stderr
    // calls os.Exit with code set as 1
}
```

### Setup (Custom Handler)

```go
    h := func(err error) {
        fmt.Println(err)
    }

    r := relay.New(h)
    defer func() { r.Filter(recover()) }()

    defer fmt.Println("reached")

    err := fail()
    r.Check(err)

    fmt.Println("should not print")}

    // Output:
    // reached
    // always fails
```

### Setup (Eased Usage)

```go
    check, filter := relay.New().Fns()
    defer func() { filter(recover()) }()

    err := fail()
    check(err)
    // prints "{cmd_name}: {err_msg}" to stderr
    // calls os.Exit with code set as 1
```

### Setup (Coded Check)

```go
    r := relay.New()
    defer func() { r.Filter(recover()) }()

    err := fail()
    r.CodedCheck(3, err)
    // prints "{cmd_name}: {err_msg}" to stderr
    // calls os.Exit with code set as first arg to r.CodedCheck
```

### Setup (Trip Function)

```go
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
```

## More Info

### Background

https://github.com/golang/go/issues/32437#issuecomment-510214015

## Documentation

View the [GoDoc](http://godoc.org/github.com/codemodus/relay)

## Benchmarks

N/A
