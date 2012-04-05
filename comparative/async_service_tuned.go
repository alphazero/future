package main

import (
	"log"
	"time"
)

// using only language primitive and type-safe requests
// (i.e. no interface{}).
func main() {

	server = StartServer()
	startClients()

	<-(make(chan int, 1))
}

type Server chan<- *request
type request struct {
	cid int
	arg time.Time
	ch  chan<- time.Duration
}

var server Server

var _SERVICE_LATENCY = int(1 * time.Nanosecond)
var _NUM_CLIENTS = 10000
var _LOAD_FACTOR = 10 * _NUM_CLIENTS
var _REPORT_LIM = 10
var _WAIT = time.Duration(1)

func startClients() {
	for i := 0; i < _NUM_CLIENTS; i++ {
		go func(cid int) {
			var n = 0
			var tocnt = 0
			var t0 time.Time = time.Now()
			var delta time.Duration
			var timeout bool

			for true {
				response := make(chan time.Duration, 1)
				request := &request{cid, time.Now(), response}
				server <- request

				// note: typically get would occur elsewhere and not immediately after request
				select {
				case <-response:
				case <-time.After(_WAIT):
					<-response
					timeout = true
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

func StartServer() Server {
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
				request.ch <- delta
			default:
				time.Sleep(1 * time.Nanosecond)
			}
		}
	}()
	return c
}
