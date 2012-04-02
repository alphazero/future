package future

import (
	"errors"
	"time"
)
/* The Untyped implementation of Future (api) */

// ----------------------------------------------------------------------------
// This structure supports future.Result interface.
//
// a basic struct that either holds a generic (interface{}) reference
// or an error reference. It is used to generically send and receive fchan results
// through channels.
//
type result_str struct {
	// v is a reference to either a result value or error value
	v interface{}

	// flag for v semantics
	faulted bool
}

// Support future.Result#Value() interface
func (r *result_str) Value() interface {} {
	if r.faulted {
		return nil
	}
	return r.v
}

// Support future.Result#Error() interface
func (r *result_str) Error() error {
	if !r.faulted {
		return nil
	}
	return r.v.(error)
}


// This structure supports future.Future interface.
type fobj_str struct {
	// channel for sending the result
	fchan chan Result
	// flag to prevent multiple sets
	finalized bool
}


// Creates a new untyped Future object to be managed by the caller (i.e. the provider)
// and returns the future.Future interface to caller.
func NewUntypedFuture() Future {
	fo := new(fobj_str)
	fo.finalized = false
	fo.fchan = make(chan Result, 1)
	return fo
}

// Future#Set support
func (f *fobj_str) SetError(e error) error {
	if f.finalized {
		// TODO: log this
		return errors.New("illegal state @ setError: already set")
	}
	f.finalized = true
	f.fchan <- &result_str{e, true}
	return nil
}

// Future#Set support
func (f *fobj_str) SetValue(v interface{}) error {
	if f.finalized {
		// TODO: log this
		return errors.New("illegal state @ setError: already set")
	}
	f.fchan <- &result_str{v, false}
	return nil
}

func (f *fobj_str) FutureResult() FutureResult {
		return (frchan)(f.fchan)
}
// used for read only channel for consumer's end
// supports future.FutureResult
type frchan <-chan Result

// support for FutureResult#Get
func (ch frchan) Get() (r Result) {
	r = <-ch
	return
}

// support for FutureResult#TryGet
func (ch frchan) TryGet(ns time.Duration) (r Result, timeout bool) {
	select {
	case r = <-ch:
		break
	case <-time.After(ns):
		timeout = true
	}
	return
}
