package future

import (
	"time"
	"bytes"
)

// basic test data
type tspec_future struct {
	data  []byte
	delay time.Duration
}

// create the test spec (test data for now)
func testspec_future() tspec_future {
	var testspec tspec_future

	// []byte data to be used
	data := "Hello there!"
	testspec.data = bytes.NewBufferString(data).Bytes()

	// using a timeout of 100 (msecs)
	delaysecs := 100
	testspec.delay = time.Duration(delaysecs) * time.Millisecond

	return testspec
}

// TODO: figure out how to plugin different impls ..
