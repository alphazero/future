package main

import (
	"log"
	"time"
)

type future <-chan interface{}

func (c future) Get() interface{} {
	return <-c
}

func (c future) TryGet(wait time.Duration) (v interface{}, timeout bool) {
	select {
	case v = <-c:
		break
	case <-time.After(wait):
		timeout = true
	}
	return
}

type pchan chan<- interface{}

func (c pchan) Set(v interface{}) {
	c <- v
}

func NewFuture() (future, pchan) {
	c := make(chan interface{}, 1)
	return future(c), pchan(c)
}

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
	fch, pch := NewFuture()

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
				ch.Set(delta)

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
			//			var result interface{}

			for true {
				// make the request and get the future.FutureResult
				future := AsyncService(cid, time.Now())

				// try get result
//				select {
//				case <-future:
//					break
//				case <-time.After(time.Duration(0)):
//					timeout = true
//					<-future
//				}

				_, timeout = future.TryGet(time.Duration(0))
				if timeout {
					future.Get()
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
