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
// Spec-Lite and Usage:
//
// Futures are designed for asynchronous hand-off of a reference between 1
// or 2 (typically 2) 'threads' of execution.  Conceptually, it is nothing
// more than wrapping a standard pattern of using a channel (of size 1) to
// transfer a reference or a value.  API is designed for minimal overhead
// and delegates any required type-safety to the programmer user.
//
// The hand-off semantics are defined by the 'Future', and 'Provider' interfaces,
// where it is expected that the instantiator of future objects is the entity
// that has exposed an API returning future.Future(s).
//
// The general pattern of use is as follows:
//
// The implementation of a function | interface-method returning future.Future
// will instantiate a type supporting the future.Future & future.Provider
// interfaces.  For example:
//
//     func AsyncSomething(...) (future.Future, ...)
//
// This function | interface-method has the role of the future.Provider.
// For example:
//
//     // here using the future package's reference implementation
//     func RemoteServiceFoo (...) (response future.Future, ...) {
//          response = future.NewUntypedFuture()
//
//          // async func invoke the service call
//          // using the future.Provider interface to fulfill the
//          // future contract
//          go func(future future.Provider) {
//              // call the remote service
//              svcresp, e := invokeRemoteService(...);
//              switch {
//              case e != nil:
//                  future.SetError(e)
//              default:
//                  future.SetValue(svcresp)
//              }
//          }(response)
//     }
//
// The consumer of the function | interface-method returning future.Future
// objects can (a) fire-and-forgetall-but-error, (b) block until results are
// provided, or (c) wait for a specified time (for cases such as meeting SLAs):
//
//       // ... call site to an async function | interface-method
//       futureResponse, ... := RemoteServiceFoo(...)
//
//       ----------------------------
//
//       // idiom (a)
//       // fire and forgetall but the error
//       go func(future future.Future, errout io.Writer ) {
//           result := future.Get()
//           if result.IsError() {
//               e := result.Error()
//               errout.Write(e.Error())
//           }
//       }(futureResponse, os.Stderr)
//
//       ----------------------------
//
//       // idiom (b)
//       // block until future results provided
//       result := futureResponse.Get()
//       if result.IsError() {
//           /* handle error */
//       } else {
//           /* process response */
//           response := result.Value().(ExpectedResponseTypeHere)
//           ...
//       }
//
//
//       ----------------------------
//
//       // idiom (c)
//       // wait for a given duration (per your SLA)
//       result, timeout := futureResponse.TryGet(slaMaxLatencyDuration)
//       if timeout {
//           /* handle failture to meet SLA case */
//           ...
//           return
//       }
//       // congrats. Got the reponse. so now just process it per (b)
//       if result.IsError() {
//           /* handle error */
//       } else {
//           /* process response */
//           response := result.Value().(ExpectedResponseTypeHere)
//           ...
//       }
//
// The behavior of future.Future is _only specified_ given the constraint that
// the receiving party adheres to the following:
//
//  a) repeated subsequent calls to FutureResult#TryGet can be made if the
//     calls result in timeouts.
//
//  b) future.Future#Get must only be called once.  It may be called in isolation
//     or can be called after one or more calls to TryGet, per timeout spec of
//     `a` above).  Any other pattern of use is unspecified and is considered
//     a programmer error.
//
// The Result interface reference obtained by the receiver per above is conceptually
// a 'union' between an 'error' or 'value' (both regardless typed as interface{}).
//
// This package makes no assumptions about the ownership and life-cycle of the
// references (whether error or value) handed off via the futures, or for that
// matter, the references to Future objects themselves.  Typically, the life-cycle
// of all objects created are close to the event bounds of creating the Future and
// obtaining either an error or value result, but that is entirely up to the
// user of the package.  The untyped implementation does not maintain any references
// to the objects it creates.
//
// Compliant implentations are required to release any rosources associated with a
// future.Future object on a successful future.Future#Get | future.Future#TryGet.
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
}
