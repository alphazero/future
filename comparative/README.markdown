# What is it?

A simple comparative use of futures and bare channels, to understand the implications in terms of both
code (readability, loc) and performance impact.

## How to try the tests

Run the future version:

   go run async_service_futures.go

Run the pure go channel version:

   go run async_service_channels.go

## code

LOCs are actually pretty much the same -- future version is shorter by 1.


## sample runs

I've tried running with various settings -- see _NUM_CLIENTS, etc. in the sources -- and it appears that there is a
fixed cost overhead, but nothing in order of mag range, for the future based version.

For the following, settings were:

	- 10000 client go routines
	- 10 nsec service latency
	- 1 nsec wait before retry on response channels.

Result in each case is dumped by a single sampled client.

Using futures:

	2012/04/05 13:03:59 (sample of 10000) : 0010 requests in 1086658000 nsec with 5 timeouts (recovered)
	2012/04/05 13:04:00 (sample of 10000) : 0010 requests in 1136708000 nsec with 6 timeouts (recovered)
	2012/04/05 13:04:01 (sample of 10000) : 0010 requests in 1119443000 nsec with 5 timeouts (recovered)
	2012/04/05 13:04:02 (sample of 10000) : 0010 requests in 1137239000 nsec with 7 timeouts (recovered)
	2012/04/05 13:04:03 (sample of 10000) : 0010 requests in 1117865000 nsec with 5 timeouts (recovered)
	2012/04/05 13:04:04 (sample of 10000) : 0010 requests in 1132970000 nsec with 5 timeouts (recovered)
	2012/04/05 13:04:05 (sample of 10000) : 0010 requests in 1170051000 nsec with 4 timeouts (recovered)
	2012/04/05 13:04:07 (sample of 10000) : 0010 requests in 1145685000 nsec with 4 timeouts (recovered)
	2012/04/05 13:04:08 (sample of 10000) : 0010 requests in 1135058000 nsec with 5 timeouts (recovered)
	2012/04/05 13:04:09 (sample of 10000) : 0010 requests in 1181092000 nsec with 3 timeouts (recovered)
	2012/04/05 13:04:10 (sample of 10000) : 0010 requests in 1153503000 nsec with 7 timeouts (recovered)

Using channels:

	2012/04/05 13:06:41 (sample of 10000) : 0010 requests in 956035000 nsec with 10 timeouts (recovered)
	2012/04/05 13:06:42 (sample of 10000) : 0010 requests in 1011511000 nsec with 10 timeouts (recovered)
	2012/04/05 13:06:43 (sample of 10000) : 0010 requests in 1006962000 nsec with 10 timeouts (recovered)
	2012/04/05 13:06:44 (sample of 10000) : 0010 requests in 996946000 nsec with 10 timeouts (recovered)
	2012/04/05 13:06:45 (sample of 10000) : 0010 requests in 1037632000 nsec with 10 timeouts (recovered)
	2012/04/05 13:06:46 (sample of 10000) : 0010 requests in 1025601000 nsec with 10 timeouts (recovered)
	2012/04/05 13:06:47 (sample of 10000) : 0010 requests in 1029911000 nsec with 10 timeouts (recovered)
	2012/04/05 13:06:48 (sample of 10000) : 0010 requests in 1013163000 nsec with 10 timeouts (recovered)
	2012/04/05 13:06:49 (sample of 10000) : 0010 requests in 1043335000 nsec with 10 timeouts (recovered)
	2012/04/05 13:06:50 (sample of 10000) : 0010 requests in 1023341000 nsec with 10 timeouts (recovered)
	2012/04/05 13:06:51 (sample of 10000) : 0010 requests in 1031367000 nsec with 10 timeouts (recovered)

Overall, both servers are capped at ~ 100k/sec request throughput.
Roughly (per above), something like 1 microsecond overhead per request for the future based version.
