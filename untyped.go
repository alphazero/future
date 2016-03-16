package future

import (
	"errors"
	"time"
)

/* The Untyped implementation of Future (api) */

// ----------------------------------------------------------------------------
// Result Value
// ----------------------------------------------------------------------------

// untyped supports ~generic future.Result interface.
//
// wraps an untyped (interface{}) value reference that
// is either the future result value or an error.
type result struct {
	v       interface{} // value ref. w/ modal semantics
	isError bool        // determines v value semantics
}

// interface: future.Result#Value()
func (r *result) Value() (v interface{}) {
	if !r.isError {
		v = r.v
	}
	return
}

// interface: future.Result#Error()
func (r *result) Error() (err error) {
	if r.isError {
		err = r.v.(error)
	}
	return
}

// interface: future.Result#Error()
func (r *result) IsError() bool {
	return r.isError
}

// ----------------------------------------------------------------------------
// Future Object
// ----------------------------------------------------------------------------

// future.futureResult supports future.Future and future.Provider
// Instances of this object are created by the future.Result provider,
// and returned to the call site as future.Future references.
type futureResult struct {
	rchan     chan Result
	finalized bool // prevent multiple sets
}

// Creates a new untyped Future object.
func NewUntypedFuture() *futureResult {
	return &futureResult{
		rchan:     make(chan Result, 1),
		finalized: false,
	}
}

// ______________________________________________________________________
// support for future.Future

// interface: future.Future#Get
func (p *futureResult) Get() (r Result) {
	r = <-(p.rchan)
	return
}

// interface: future.Future#TryGet
func (p *futureResult) TryGet(ns time.Duration) (r Result, timeout bool) {
	select {
	case r = <-(p.rchan):
	case <-time.After(ns):
		timeout = true
	}
	return
}

// ______________________________________________________________________
// support for future.Provider

func (f *futureResult) SetError(e error) error {
	if f.finalized {
		return errors.New("illegal state @ setError: already set")
	}
	f.set(&result{e, true})
	return nil
}

func (f *futureResult) SetValue(v interface{}) error {
	if f.finalized {
		return errors.New("illegal state @ setValue: already set")
	}
	f.set(&result{v, false})
	return nil
}

func (f *futureResult) set(r Result) {
	f.rchan <- r
	f.finalized = true
	close(f.rchan)
}
