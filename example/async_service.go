package main

import (
	"future"
	"log"
	"time"
)

// usage example
// creating and using an asynchronous service
func main() {

	server = startServer()
	startClients()

	var never chan struct{}
	<-never
}

const (
	latencyFactor = 1               // latency of each service request
	numClients    = 10              // number of concurrent clients
	loadFactor    = 10 * numClients // load factor
)

type asyncRequest struct {
	cid  uint
	t0   time.Time
	fobj future.Future
}

var server chan<- *asyncRequest

// --- the clients -----------------------------------

func startClients() {
	for i := uint(0); i < numClients; i++ {
		go func(cid uint) {
			log.Printf("client %d started\n", cid)
			var n = 0
			var tocnt = 0
			var t0 time.Time = time.Now()

			for {
				// make the request and get the future.FutureResult
				// note: calling time.Now in loop significantly impacts
				//       performance. sampled results values do not reflect
				//       actual future usage perf. cost.
				fresult := callService(cid, time.Now())

				// TryGet the future result
				// note: tryget & then get on timeout is non-optimal
				//       used here only to show api usage
				result, timeout := fresult.TryGet(time.Microsecond * 10)
				if timeout {
					result = fresult.Get() // wait for it
				}

				// client 0 will dump its results as a sample
				if cid == 0 {
					switch {
					case timeout:
						tocnt++
					case result.IsError(): /* nop - just api demo */
					}
					n++
					if n == 1000 {
						delta := time.Now().Sub(t0)
						log.Printf("(sample of %d) : %04d requests in %d nsec with %d timeouts\n", numClients, n, delta, tocnt)
						n = 0
						tocnt = 0
						t0 = time.Now()
					}
				}
			}
		}(i)
	}
}

// --- the service -----------------------------------

func callService(cid uint, t0 time.Time) future.FutureResult {

	// create the Future object
	fobj := future.NewUntypedFuture()

	// create an asyncRequest
	// server will use future object to post its response
	request := &asyncRequest{cid, time.Now(), fobj}

	// queue the request
	server <- request

	// Return the FutureResult of the Future object to the caller
	return fobj.Result()
}

func startServer() chan<- *asyncRequest {
	c := make(chan *asyncRequest)
	go func() {
		for {
			// process request queue or sleep if none pending
			//
			select {
			case request := <-c:
				// sleep for fake service latency
				time.Sleep(time.Duration(latencyFactor) * time.Nanosecond)

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
