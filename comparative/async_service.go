package main

import (
	"future"
	"log"
	"time"
)

// using UntypedFuture -- request is type-safe;
// response is interface{}
func main() {

	server = StartServer()
	startClients()

	<-(make(chan int, 1))
}

type Server chan<- *request
type request struct {
	cid    int
	arg    time.Time
	future future.Future
}

var server Server

var _SERVICE_LATENCY = int(10 * time.Nanosecond)
var _NUM_CLIENTS = 10000
var _LOAD_FACTOR = 10 * _NUM_CLIENTS
var _REPORT_LIM = 10
var _WAIT = time.Duration(0)

func startClients() {
	for i := 0; i < _NUM_CLIENTS; i++ {
		go func(cid int) {
			var n = 0
			var tocnt = 0
			var t0 time.Time = time.Now()
			var delta time.Duration
			var timeout bool

			for true {
				future := future.NewUntypedFuture()
				request := &request{cid, time.Now(), future}
				server <- request

				// note: typically get would occur elsewhere and not immediately after request
				_, timeout = future.Result().TryGet(time.Duration(_WAIT))
				if timeout {
					future.Result().Get()
				}

				// client 0 will dump its results as a sample
				if cid == 0 {
					if timeout {
						tocnt++
					}

					n++
					if n >= _REPORT_LIM {
						delta = time.Now().Sub(t0)
						log.Printf("(sample of %d) : %04d requests in %d nsec with %d timeouts (recovered)\n", _NUM_CLIENTS, n, delta, tocnt)
						n = 0
						tocnt = 0
					}
					if n == 0 {
						t0 = time.Now()
					}
				}
			}
		}(i)
	}
}

func StartServer() chan<- *request {
	c := make(chan *request)
	go func() {
		for {
			// process request queue or sleep if none pending
			select {
			case request := <-c:
				// sleep for fake service latency
				time.Sleep(time.Duration(_SERVICE_LATENCY) * time.Nanosecond)

				// result is just the delta of t0 of request and time now
				t0 := request.arg
				delta := time.Now().Sub(t0)
				request.future.SetValue(delta)
			default:
				time.Sleep(1 * time.Nanosecond)
			}
		}
	}()
	return c
}
