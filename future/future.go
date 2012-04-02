// the future package define the semantics of asynchronous (future value)
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
// Futures
// result and future result channel

// Result defines a basic struct that either holds a generic (interface{}) reference
// or an error reference. It is used to generically send and receive future results
// through channels.

type Result struct {
	// Value of the Result.
	// Value may be nil.
	Value interface{}

	// Error status of the future Result.
	// if Error is NOT nil, then Value must be disregarded
	// and Result considered to be an error
	Error error
}

// A future defines the interface to a Result that is provided asynchronously
// per the standard practice semantics of a 'future'
//
// REVU(jh): split into FutureProvider and Future
type Future interface {
	// sets the value of the future Result
	// if future Result was an error, v should be nil
	// if e is nil, future Result is considered to be valid, even if v is nil
	//
	// REVU: should be for future provider
	Set(v interface{}, e error)

	// Blocks until the future Result is provided
	// REVU: should be future consumer
	Get() (r *Result)

	// Blocking call with timeout.
	// Returns timeout value of TRUE if ns Duration elapsed; r will be nil.
	// Returns r value if value was provided before ns timeout duration elapsed.
	// REVU: should be future consumer
	TryGet(ns time.Duration) (r *Result, timeout bool)
}

// ----------------------------------------------------------------------------
// Core support

// future type is a channel of future results
// supporting the Future interface
type future chan *Result

// Create a new future
func NewFuture() future {
	return make(future, 1)
}

// Future#Set support
func (f future) Set(v interface{}, e error) {
	f <- &Result{v, e}
}

// Future#Get support
func (f future) Get() (r *Result) {
	r = <-f
	return
}

// Future#TryGet support
func (f future) TryGet(ns time.Duration) (r *Result, timeout bool) {
	select {
	case r = <-f:
		break
	case <-time.After(ns):
		timeout = true
	}
	return
}

