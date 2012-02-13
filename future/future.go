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

// Future
//
// A future defines the interface to a Result that is provided asynchronously
// per the standard practice semantics of a 'future'
type Future interface {
	// sets the value of the future Result
	// if future Result was an error, v should be nil
	// if e is nil, future Result is considered to be valid, even if v is nil
	Set(v interface{}, e error)

	// Blocks until the future Result is provided
	//
	Get() (r Result)
	TryGet(ns time.Duration)
}

// ----------------------------------------------------------------------------
// Support

// future type is a channel of FutureResults
// supporting the Future interface
type future chan Result

// Create a new future - channel depth is 1
func NewFuture() future {
	return make(future, 1)
}

// Future#Set support
func (f future) Set(v interface{}, e error) {
	f <- Result{v, e}
}

// Future#Get support
func (f future) Get() (r Result) {
	r = <-f
	return
}

// Future#TryGet support
func (f future) TryGet(ns time.Duration) (r Result, timeout bool) {
	select {
	case r = <-f:
		break
	case <-time.After(ns):
		timeout = true
	}
	return
}

/* 
// idea about generics that came to mind based on above:
// 
// proposal for a new 'generic' special type interface{*} used
// only in type definitions -- for now limited to func type defs.
//
// use of this type allows for compile-time matching at call site
// Type cast errors remain runtime as with existing Go mechanisms
//
// example usage

// v is the std. Go interface{} and nothing is changed here
// *func is additional compile time syntax to prevent programmer error
// returned result interface{*} is compile time syntax indicating generic type
//
type Converter func? (v interface{}) interface{?}

// Runtime issues remain as they are with current Go.
// But we hugely reduce syntactic noise of breaking KISS
// with boiler plate type-safe implementations.
//
// ex: this function should be compile time mappable to Converter type.
//
func boolConverter (v interface{}) bool {
	return v.(bool)
}

// and this statement should compile
//
var BooleanConverter Converter = boolConverter

// ex continued --
//      // call site compile-time check insures
//      // receiver 'flag' is of correct expected type
//
//      var flag bool = BooleanConverter(foo)

// extending idea for finer compile time control
//
type Converter func? (v interface{}) interface{specific-interface-type}

// extending idea for structs
//
// define a generic struct
type Result?T struct {
	Value  interface{T}
	Error  error
}

// we could then state
//
type BooleanResult Result?bool

// which would map to
type Boolean struct {
  Value bool
  Error error
}

type Boolean Generic?bool


*/
