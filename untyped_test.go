/*
 * Tests the untyped implementation of the Future API.
 */
package future

import (
	"bytes"
	"fmt"
	"testing"
)

// test Future - all aspect, expect blocking Get
// are tested.  A future object is created (but not set)
// and timeout and error tests are made.  Then value is set
// and test repeated.
func TestFutureContract(t *testing.T) {

	// prep data
	tspec := testspec_future()

	// using untyped future
	future := NewUntypedFuture()
	fr := future.FutureResult()

	/* TEST timed call to uninitialized (not Set) future value,
	 * expecting a timeout.
	 * MUST return timeout of true
	 * MUST return value of nil
	 */
	fvalue1, timedout1 := fr.TryGet(tspec.delay)
	if !timedout1 {
		t.Error("BUG: timeout expected")
	}
	if fvalue1 != nil {
		t.Error("Bug: value returned: %s", fvalue1)
	}

	// set the future result
	future.SetValue(tspec.data)

	/* TEST timed call to initialized (set) future value
	 * expecting data, no error and no timeout
	 * MUST return timeout of false
	 * MUST return error of nil
	 * MUST return value equal to data
	 */
	fvalue2, timedout2 := fr.TryGet(tspec.delay)
	if timedout2 {
		t.Error("BUG: should not timeout")
	}
	if fvalue2 == nil {
		t.Error("Bug: should not return future nil")
	} else {
		if fvalue2.Error() != nil {

		}
		if fvalue2.Value() == nil {

		}
		value := fvalue2.Value().([]byte)
		if bytes.Compare(value, tspec.data) != 0 {
			t.Error("Bug: future value not equal to data set")
		}
	}
}

func TestFutureWithBlockingGet(t *testing.T) {

	// prep data
	// prep data
	tspec := testspec_future()

	// using basic FutureBytes
	future := NewUntypedFuture()
	fr := future.FutureResult()

	// test go routine will block on Get until
	// value is set.
	sig := make(chan bool, 1)
	go func() {
		/* TEST timed call to initialized (set) future value
		 * expecting data, no error and no timeout
		 * MUST return error of nil
		 * MUST return value equal to data
		 */
		fvalue := fr.Get()
		if fvalue == nil {
			t.Error("Bug: should not return future nil")
		} else {
			if fvalue.Error() != nil {
				t.Error("Bug: unexpected error %s", fvalue.Error())
			}
			value := fvalue.Value().([]byte)
			if bytes.Compare(value, tspec.data) != 0 {
				t.Error("Bug: future value not equal to data set")
			}
		}
		sig <- true
	}()

	// set the data
	future.SetValue(tspec.data)

	<-sig

}

func TestFutureTimedBlockingGet(t *testing.T) {
	// tests timed blocking gets with no errors

	// prep data
	// prep data
	tspec := testspec_future()

	// using basic FutureBytes
	future := NewUntypedFuture()
	fr := future.FutureResult()

	// test go routine will block on Get until
	// value is set or timeout expires
	sig := make(chan bool, 1)
	go func() {
		/* TEST timed call to initialized (set) future value
		 * expecting data, no error and no timeout
		 * MUST return error of nil
		 * MUST return value equal to data
		 */
		fvalue, timedout := fr.TryGet(tspec.delay)
		if timedout {
			t.Error("BUG: should not timeout")
		}
		if fvalue == nil {
			t.Error("Bug: should not return future nil")
		} else {
			if fvalue.Error() != nil {
				t.Error("Bug: unexpected error %s", fvalue.Error())
			}
			if fvalue.Value() == nil {
				t.Error("Bug: value is nil")
			} else {
				value := fvalue.Value().([]byte)
				if bytes.Compare(value, tspec.data) != 0 {
					t.Error("Bug: future value not equal to data set")
				}
			}
		}
		sig <- true
	}()

	// set the data
	future.SetValue(tspec.data)

	<-sig

}

// dummy error type used for testing
type FoobarError int64

func (e FoobarError) Error() string {
	return fmt.Sprintf("%d", e)
}

func TestFutureTimedBlockingGetWithError(t *testing.T) {
	// tests timed blocking gets with no errors

	// prep data
	// prep data
	tspec := testspec_future()

	// using basic FutureBytes
	future := NewUntypedFuture()
	fr := future.FutureResult()

	// test go routine will block on Get until
	// value is set or timeout expires
	sig := make(chan bool, 1)

	// error code we will set on future value
	var errorCode FoobarError = 111
	go func() {
		/* TEST timed call to initialized (set) future value
		 * expecting data, no error and no timeout
		 * MUST return error of nil
		 * MUST return value equal to data
		 */
		fvalue, timedout := fr.TryGet(tspec.delay)
		if timedout {
			t.Error("BUG: should not timeout")
		}
		if fvalue == nil {
			t.Error("Bug: should not return future nil")
		} else {
			error := fvalue.Error()
			if error == nil {
				t.Error("Bug: expected error")
			} else {
				value := fvalue.Value()
				if value != nil {
					t.Error("Bug: future value must be nil if error is set.")
				}
				if error != errorCode {
					t.Error("Bug: expected error code of ", errorCode)
				}
			}
		}
		sig <- true
	}()

	// set the data
	// note we are setting a future result with error.
	var e FoobarError = FoobarError(111) // an error
	future.SetError(e)

	<-sig
}
