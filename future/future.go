// the fchan package define the semantics of asynchronous (fchan value)
// data access.
//
// The API provided can be used as either generic or type-safe
// manner (though the utility of the former case is limited to wrapping a
// well-known channel/timer select.)
//
// TODO:
// The package is to provide type-safe variants for builtin types.  Users
// can (and are expected) to write their own variants for custom types (until
// Go gets generics.)
//
// API design was strongly inspired by Java's Futures.
package future

import (
	"time"
)

// ----------------------------------------------------------------------------
// Futures API
// ----------------------------------------------------------------------------

type Result interface {
	// Value result.
	// Value of the Result - Value may be nil ONLY in case of Errors.
	Value() interface{}

	// Error result.
	// if Error is NOT nil, then Value must be nil.
	Error() error
}

// The consumer end of the Future that is provided to the end-user.
// This interface defines the semantics of asynchronous access to the
// future result.
//
type FutureResult interface {
	// Blocks until the fchan Result is provided
	Get() Result

	// Blocking call with timeout.
	// Returns timeout value of TRUE if ns Duration elapsed; r will be nil.
	// Returns r value if value was provided before ns timeout duration elapsed.
	TryGet(ns time.Duration) (r Result, timeout bool)
}

type Future interface {

	// sets the value of the fchan Result
	// Future.Value will be nil
	// A non-nil error is returned if already set.
	SetError(e error) error

	// sets an erro fchan Result - note that nil values are NOT permitted.
	// Future.Error will be nil
	// A non-nil error is returned if already set.
	SetValue(v interface{}) error

	// FutureResult is provided to the consumer of the future result.
	FutureResult() FutureResult
}
