// The future package define the semantics of asynchronous (future value)
// data access.
//
// The API can be used as either generic or type-safe manner:
// the utility of the former case is limited to wrapping a
// well-known channel/timer select pattern in a standard, tested, and
// modularized manner; and the latter obviates the need for casts but
// they do remain runtime mechanisms due to lack of support for generics.
//
// The package may provide type-safe variants for builtin types.  Users
// can (and are free) to write their own variants for custom types, if
// that is a requirement.  A reference implementation supporting untyped
// futures is provided, and can be used via exported future.NewUntypedFuture.
//
// API design was strongly inspired by Java's Futures.
//
// -
//
//  Spec-Lite and Usage:
//
//  Futures are designed for asynchronous hand-off of a reference between 1
//  or 2 (typically 2) 'threads' of execution.  Conceptually, it is nothing
//  more than wrapping a standard pattern of using a channel (of size 1) to
//  transfer a reference or a value.  API is designed for minimal overhead
//  and delegates any required type-safety to the programmer user.
//
//  The hand-off semantics are defined by the Future interface, which is
//  expected by design to be 'owned' by the party making promises of delivery
//  'in future' to the other party.
//
//  The Future object provides the basic check of enforcing a one time only
//  setting of the future results (to either a value or an error reference).
//  But it does (and can) not prevent other expected usage constraints.
//
//  The holder of the reference ot the Future object is expected to not leak the
//  reference to the Future
//
//  The 'receiver' of the data interacts in the hand-off via the FutureResult
//  interface.  A reference to this interface can be obtained from the Future
//  interface (above).  This party can use either the blocking or non-blocking
//  with timeout functions of FutureResult interface to obtain the value (or
//  a application specific error object minimally supporting Go builtin Error
//  interface.
//
//  The blocking FutureResult#Get will block until the promised results are
//  provided.  As of now, it is unspecified if this function should support
//  interrupts or not.
//
//  The non-blocking FutureResult#TryGet(timeout in nano-seconds) must return
//  no later than after the timeout duration has elapsed after making the call,
//  in the idealized case.  It may return before the timeout duration has elapsed
//  with a result.
//
//  The behavior of FutureResult is only specified given the constraint that
//  the receiving party adheres to the following:
//
//  a) repeated subsequent calls to FutureResult#TryGet can be made if the
//     calls result in timeouts.
//
//  b) FutureResult#Get must only be called once.  It may be called in isolation
//     or can be called after one or more calls to TryGet, per timeout spec of
//     `a` above).  Any other pattern of use is unspecified.
//
//
//  The Result interface reference obtained by the receiver per above is conceptually
//  a 'union' between an 'error' or 'value' (both regardless typed as interface{}).
//
//  This package makes no assumptions about the ownership and life-cycle of the
//  references (whether error or value) handed off via the futures, or for that
//  matter, the references to Future objects themselves.  Typically, the life-cycle
//  of all objects created are close to the event bounds of creating the Future and
//  obtaining either an error or value result, but that is entirely up to the
//  user of the package.  The untyped implementation does not maintain any references
//  to the objects it creates.
//
//  -
//
//  Usage in the minimal sense is quite trivial and entirely in line with the
//  idiomatic usage of goroutines and channels.  Typically, one goroutine is
//  handing off read only references to (size 1) channels which they are expected
//  to dequeue via a select construct.  Using Future, the reference handed off is
//  a reference to the FutureResult instead of a bare channel, with enhanced
//  (and extensible) semantics.  And instead of adding the result directly to
//  a channel, using futures the provider does via the Future interface api.
//
package future

import (
	"time"
)

// ----------------------------------------------------------------------------
// Future
// ----------------------------------------------------------------------------

// future.Future defines the api of future objects. Objects supporting this
// interface are created by the async function/service and returned to call site.
type Future interface {
	// Blocking get waits until future.Result is available.
	Get() (r Result)

	// Returns result or timeouts after specified wait duration
	TryGet(wait time.Duration) (r Result, timeout bool)
}

// ----------------------------------------------------------------------------
// Future Result
// ----------------------------------------------------------------------------

// future.Result defines the type-generic future value results returned via a
// future.Future.
type Result interface {
	// Value result.
	// Value of the Result - Value may be nil ONLY in case of Errors.
	Value() interface{}

	// Error result.
	// if Error is NOT nil, then Value must be nil.
	Error() error

	// convenience method
	IsError() bool
}

// ----------------------------------------------------------------------------
// Future Provider
// ----------------------------------------------------------------------------

// future.Provider defines the api for use by the provider of future.Results
type Provider interface {

	// sets the value of the fchan Result
	// Future.Value will be nil
	// A non-nil error is returned if already set.
	SetError(e error) error

	// sets an erro fchan Result - note that nil values are NOT permitted.
	// Future.Error will be nil
	// A non-nil error is returned if already set.
	SetValue(v interface{}) error

	// Result is provided to the consumer of the future result.
	// REVU: why here?
	Result() Result
}
