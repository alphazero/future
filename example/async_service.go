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

	for i := 0; i < _NUM_CLIENTS; i++ {
		go func() {
			log.Println("client started ..")
			for true {
				fresult := AsyncService(time.Now())
				//		        time.Sleep(10 * time.Nanosecond)

				// try
				result, timeout := fresult.TryGet(time.Duration(_LOAD_FACTOR * _SERVICE_LATENCY))
				if timeout {
					log.Printf("timeout\n")
				} else {
					if !result.IsError() {
						//			            delta := result.Value().(time.Duration)
						//			            log.Printf("%s\n", delta)
					}
				}
			}
		}()
	}
	flatch := future.NewUntypedFuture()
	flatch.FutureResult().Get()
}

func AsyncService(t0 time.Time) future.FutureResult {

	// create the Future object
	fobj := future.NewUntypedFuture()

	// create an asyncRequest
	// server will use future object to post its response
	request := &asyncRequest{t0, fobj}

	// queue the request
	server <- request

	// Return the FutureResult of the Future object to the caller
	return fobj.FutureResult()
}

var _SERVICE_LATENCY = 10
var _NUM_CLIENTS = 100
var _LOAD_FACTOR = 100 * _NUM_CLIENTS

type asyncRequest struct {
	t0   time.Time
	fobj future.Future
}

var server chan<- *asyncRequest

func StartServer() chan<- *asyncRequest {
	c := make(chan *asyncRequest)
	go func() {
		for {
			request := <-c
			fobj := request.fobj
			t0 := request.t0

			time.Sleep(time.Duration(_SERVICE_LATENCY) * time.Nanosecond)
			delta := time.Now().Sub(t0)
			fobj.SetValue(delta)
		}
	}()
	return c
}
