package main

import (
	"future"
	"log"
	"time"
)

// usage example -- creating and using an asynchronous service
//
func main() {

	server = StartServer()
	startClients()

	flatch := future.NewUntypedFuture()
	flatch.FutureResult().Get() // never returns ..
}

// latency of each service request
var _SERVICE_LATENCY = int(time.Microsecond)

// number of concurrent clients
var _NUM_CLIENTS = 10

// load factor
var _LOAD_FACTOR = 10 * _NUM_CLIENTS

type asyncRequest struct {
	cid  int
	t0   time.Time
	fobj future.Future
}

var server chan<- *asyncRequest

// --- the clients -----------------------------------

func startClients() {
	for i := 0; i < _NUM_CLIENTS; i++ {
		go func(cid int) {
			log.Printf("client %d started\n", cid)
			var n = 0
			var tocnt = 0
			var t0 time.Time = time.Now()
			var delta time.Duration
			var timeout bool
			var result future.Result

			for true {
				// make the request and get the future.FutureResult
				fresult := AsyncService(cid, time.Now())

				// TryGet the future result
				result, timeout = fresult.TryGet(time.Duration(0))
				if timeout {
					result = fresult.Get()
				}

				// client 0 will dump its results as a sample
				if cid == 0 {
					if timeout {
						tocnt++
					} else {
						if !result.IsError() {
							/* nop -- have the result */
						}
					}
					n++
					if n >= 1000 {
						delta = time.Now().Sub(t0)
						log.Printf("(sample of %d) : %04d requests in %d nsec with %d timeouts (recovered)\n", _NUM_CLIENTS, n, delta, tocnt)
						n = 0
						tocnt = 0
					}
					if n == 0 {
						t0 = time.Now()
					}
				}

				// sleep for 1 ns to allow server to catch up
				//				time.Sleep(1)
			}
		}(i)
	}
}

// --- the service -----------------------------------

func AsyncService(cid int, t0 time.Time) future.FutureResult {

	// create the Future object
	fobj := future.NewUntypedFuture()

	// create an asyncRequest
	// server will use future object to post its response
	request := &asyncRequest{cid, time.Now(), fobj}

	// queue the request
	server <- request

	// Return the FutureResult of the Future object to the caller
	return fobj.FutureResult()
}

func StartServer() chan<- *asyncRequest {
	c := make(chan *asyncRequest)
	go func() {
		for {
			// process request queue or sleep if none pending
			//
			select {
			case request := <-c:
				// sleep for fake service latency
				time.Sleep(time.Duration(_SERVICE_LATENCY) * time.Nanosecond)

				// get the future object from the request
				fobj := request.fobj

				// result is just the delta of t0 of request and time now
				t0 := request.t0
				delta := time.Now().Sub(t0)
				fobj.SetValue(delta)

			default:
				time.Sleep(1 * time.Nanosecond)
			}
		}
	}()
	return c
}
