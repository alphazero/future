/* white box tests */

package future

import (
	"bytes"
	"fmt"
	"testing"
	"time"
)

/// basic test spec /////////////////////////////////////////////////////

// basic test data
type testspec struct {
	data          []byte
	wait          time.Duration
	providerDelay time.Duration
	err           error
}

// create the test spec (test data for now)
func testSpec() testspec {
	data := "Salaam!"
	var testspec = testspec{
		data:          []byte(data),
		err:           fmt.Errorf("Relax. To err is Human"),
		wait:          time.Millisecond,
		providerDelay: time.Microsecond,
	}

	return testspec
}

/// tests //////////////////////////////////////////////////////////////

// ____________________________________________________________________
// general contract

// timed call to uninitialized (not Set) future value,
// expecting a timeout.
// MUST return timeout of true
// MUST return value of nil
func TestFutureContractNew(t *testing.T) {
	futureObj := NewUntypedFuture()

	// note: wait value is irrelevant in this test
	// any value is fine
	result, timeout := futureObj.TryGet(time.Microsecond)
	switch {
	case timeout == false:
		t.Error("exected timeout => true")
	case result != nil:
		t.Error("expected nil result on timeout")
	}
}

// timed call to initialized (set) future value
// expecting data, no error and no timeout
// MUST return timeout of false
// MUST return error of nil
// MUST return value equal to spec data
func TestFutureSetNonErrorValueThenTryGet(t *testing.T) {
	// test spec & data
	test := testSpec()

	// create & set the future result
	futureObj := NewUntypedFuture()
	futureObj.SetValue(test.data)

	// note: test.wait value is irrelevant in this test
	// any value is fine
	result, timeout := futureObj.TryGet(test.wait)
	switch {
	case timeout:
		t.Error("exected timeout => false")
	case result == nil:
		t.Error("expected non-nil result with timeout == false")
	case result.IsError():
		t.Error("expected IsError => false")
	case result.Error() != nil:
		t.Error("exected Error() => nil")
	case result.Value() == nil:
		t.Error("expected non-nil Value()")
	default:
		typedValue := result.Value().([]byte)
		if bytes.Compare(typedValue, test.data) != 0 {
			t.Error("unexpected result value")
		}
	}
}

// timed call to initialized (set) future value
// expecting error and no timeout
// MUST return timeout of false
// MUST return error result
// MUST return error value equal to spec error
func TestFutureSetErrorValueThenTryGet(t *testing.T) {
	// test spec & data
	test := testSpec()

	// create & set the future result
	futureObj := NewUntypedFuture()
	futureObj.SetError(test.err)

	// note: test.wait value is irrelevant in this test
	// any value is fine
	result, timeout := futureObj.TryGet(test.wait)
	switch {
	case timeout:
		t.Error("exected timeout => false")
	case result == nil:
		t.Error("expected non-nil result with timeout == false")
	case result.IsError() == false:
		t.Error("expected IsError => true")
	case result.Error() == nil:
		t.Error("exected non-nil Error()")
	case result.Value() != nil:
		t.Error("expected Value() == nil")
	default:
		if result.Error() != test.err {
			t.Error("unexpected result error value")
		}
	}
}

// ____________________________________________________________________
// ops

// blocking Get with initialized (set) future value
// expecting data, no error and no timeout
// MUST NOT timeout
// MUST return error of nil
// MUST return value equal to data
//
func TestFutureGetDelayThenSet(t *testing.T) {

	test := testSpec()

	futureObj := NewUntypedFuture()

	// test go routine will block on Get until
	// value is set. Result will be sent on rch
	rch := make(chan Result)
	go func(future Future) {
		rch <- future.Get()
	}(futureObj)

	// sleep a bit and then set the data
	wait := time.Microsecond
	time.Sleep(time.Microsecond)
	// note: explict cast not required
	// being explicit to clarify semantics
	Provider(futureObj).SetValue(test.data)

	// MUST NOT timeout
	var result Result
	select {
	case <-time.After(wait * 100):
		t.Fatalf("expected result for Get by now")
	case result = <-rch:
		/* checked below */
	}

	// MUST return no-error result
	// MUST return value equal to data
	switch {
	case result.IsError():
		t.Error("expected IsError => false")
	case result.Error() != nil:
		t.Error("exected Error() => nil")
	case result.Value() == nil:
		t.Error("expected non-nil Value()")
	default:
		typedValue := result.Value().([]byte)
		if bytes.Compare(typedValue, test.data) != 0 {
			t.Error("unexpected result value")
		}
	}
}

// blocking Get with initialized (set) future value
// expecting data, no error and no timeout
// MUST NOT timeout
// MUST return error of nil
// MUST return value equal to data
//
func TestFutureWaitGetThenSet(t *testing.T) {

	test := testSpec()

	futureObj := NewUntypedFuture()

	// test go routine will block on Get until
	// value is set. Result will be sent on rch
	rch := make(chan Result)
	go func(future Future) {
		rch <- future.Get()
	}(futureObj)

	// delay a bit and then set the data
	time.Sleep(test.providerDelay)
	// note: explict cast not required
	// being explicit to clarify semantics
	Provider(futureObj).SetValue(test.data)

	// MUST NOT timeout
	var result Result
	select {
	case <-time.After(test.providerDelay * 100):
		t.Fatalf("expected result for Get by now")
	case result = <-rch:
		/* checked below */
	}

	// MUST return no-error result
	// MUST return value equal to data
	switch {
	case result.IsError():
		t.Error("expected IsError => false")
	case result.Error() != nil:
		t.Error("exected Error() => nil")
	case result.Value() == nil:
		t.Error("expected non-nil Value()")
	default:
		typedValue := result.Value().([]byte)
		if bytes.Compare(typedValue, test.data) != 0 {
			t.Error("unexpected result value")
		}
	}
}

// blocking Get with initialized (set) future value
// expecting data, no error and no timeout
// MUST NOT timeout
// MUST return error of nil
// MUST return value equal to data
//
func TestFutureTryGetDelayThenSetBeforeTimeout(t *testing.T) {

	test := testSpec()

	futureObj := NewUntypedFuture()

	// test go routine will block on Get until
	// value is set. Result will be sent on rch
	rch := make(chan Result)
	go func(future Future) {
		var result Result
		result, timeout := future.TryGet(test.wait)
		if timeout {
			t.Fatalf("TryGet timeout")
		}
		rch <- result
	}(futureObj)

	// sleep a bit and then set the data
	time.Sleep(test.providerDelay)
	// note: explict cast not required
	// being explicit to clarify semantics
	Provider(futureObj).SetValue(test.data)

	// MUST NOT timeout
	var result Result
	select {
	case <-time.After(test.providerDelay * 100):
		t.Fatalf("expected result for Get by now")
	case result = <-rch:
		/* checked below */
	}

	// MUST return no-error result
	// MUST return value equal to data
	switch {
	case result.IsError():
		t.Error("expected IsError => false")
	case result.Error() != nil:
		t.Error("exected Error() => nil")
	case result.Value() == nil:
		t.Error("expected non-nil Value()")
	default:
		typedValue := result.Value().([]byte)
		if bytes.Compare(typedValue, test.data) != 0 {
			t.Error("unexpected result value")
		}
	}
}
