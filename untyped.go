package future

import (
	"errors"
)

/* The Untyped implementation of Future (api) */

// ----------------------------------------------------------------------------
// Result Value
// ----------------------------------------------------------------------------

// untyped supports ~generic future.Result interface.
//
// a basic struct that either holds a generic (interface{}) reference
// or an error reference. It is used to generically send and receive fchan results
// through channels.
//
type result struct {
	// v is a reference to either a result value or error value
	v interface{}

	// flag for v semantics
	faulted bool
}

// future.Result#Value() interface
func (r *result) Value() interface{} {
	if r.faulted {
		return nil
	}
	return r.v
}

// future.Result#Error() interface
func (r *result) Error() error {
	if !r.faulted {
		return nil
	}
	return r.v.(error)
}

// future.Result#Error() interface
func (r *result) IsError() bool {
	return r.faulted
}

// ----------------------------------------------------------------------------
// Future Object handle
// ----------------------------------------------------------------------------

// supports future.Future interface.
type fobj_str struct {
	fchan     chan Result // channel for sending the result
	finalized bool        // flag to prevent multiple sets
}

// Future#Set support
func (f *fobj_str) SetError(e error) error {
	if f.finalized {
		return errors.New("illegal state @ setError: already set")
	}
	f.finalized = true
	f.fchan <- &result{e, true}
	return nil
}

// Future#Set support
func (f *fobj_str) SetValue(v interface{}) error {
	if f.finalized {
		return errors.New("illegal state @ setError: already set")
	}
	f.fchan <- &result{v, false}
	return nil
}

func (f *fobj_str) Result() FutureResult {
	return f.fchan
}

// ----------------------------------------------------------------------------
// Functions
// ----------------------------------------------------------------------------

// (exp.)
var UntypedBuilder Builder = NewUntypedFuture

// Creates a new untyped Future object.
func NewUntypedFuture() Future {
	fo := new(fobj_str)
	fo.finalized = false
	fo.fchan = make(chan Result, 1)
	return fo
}
