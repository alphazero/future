package main

import (
	"log"
	"time"
)

type future <-chan time.Duration
type pchan chan<- time.Duration

type Server chan<- *request

var server Server

type request struct {
	cid int
	arg time.Time
	ch  pchan
}

// latency of each service request
var _SERVICE_LATENCY = int(time.Microsecond)

// number of concurrent clients
var _NUM_CLIENTS = 10

// load factor
var _LOAD_FACTOR = 10 * _NUM_CLIENTS

func AsyncService(cid int, t0 time.Time) future {

	// create the Future chan
	c := make(chan time.Duration, 1)
	fch := future(c)
	pch := pchan(c)

	// create an asyncRequest
	// server will use future object to post its response
	request := &request{cid, time.Now(), pch}

	// queue the request
	server <- request

	// Return the FutureResult of the Future object to the caller
	return fch
}

func StartServer() Server {
	c := make(chan *request)
	go func() {
		for {
			// process request queue or sleep if none pending
			//
			select {
			case request := <-c:
				// sleep for fake service latency
				time.Sleep(time.Duration(_SERVICE_LATENCY) * time.Nanosecond)

				// get the future object from the request
				ch := request.ch

				// result is just the delta of t0 of request and time now
				t0 := request.arg
				delta := time.Now().Sub(t0)
				ch <- delta

			default:
				time.Sleep(1 * time.Nanosecond)
			}
		}
	}()
	return c
}

func startClients() {
	for i := 0; i < _NUM_CLIENTS; i++ {
		go func(cid int) {
			log.Printf("client %d started\n", cid)
			var n = 0
			var tocnt = 0
			var t0 time.Time = time.Now()
			var delta time.Duration
			var timeout bool
			//			var result time.Duration

			for true {
				// make the request and get the future.FutureResult
				future := AsyncService(cid, time.Now())

				// try get result

				select {
				case <-future:
					break
				case <-time.After(time.Duration(0)):
					timeout = true
					<-future // block for it
				}

				// client 0 will dump its results as a sample
				if cid == 0 {
					if timeout {
						tocnt++
					} else {
						/* nop -- have the result */
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

func main() {

	server = StartServer()
	startClients()

	<-(make(chan int, 1))
}
