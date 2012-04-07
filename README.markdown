# package future

Package future provides for asynchronous data access per the future pattern.  It is a minimal treatment over the
existing Go channels that allows the user to opt for partial or full use of the provided features.

## background and motivation

Go is a very interesting language for various reasons.  One very interesting aspect of Go, in my opinion, is that it
resists misguided complexity.  It is also an excellent bottom-up design language, in the sense that it is quite natural
to code in the small and then incrementally add features as necessary.

I wrote the first version of this, when Go came out a few years ago, in Go-Redis asynchronous pipelines, and the intent
was to extract it as a distinct package.  Naturally it was basically a rendition of Java's Futures and how they were
used in JRedis (the Java version of Go-Redis).  And then for various reasons stopped writing more Go code.

Fast forward a few years to 2012.  I started coding in Go again and in course of relearning the language started with
the previously planned future extraction.  The main motivation at this point is to (learn to) express design using
idiomatic Go, and towards that learning goal this minor exercise has already been helpful, personally.

## fine we're happy for you, but is it of any use to the rest of us?

Possibly.

The /comparative directory of the project contains two (effectively) identical programs.  Both in terms of LOC, and
performance, it is quite close.  I personally find the future based code more clear in terms of intent, but that is
a subjective view. And while a bit of selector code pretty much gives you basic 'futures' in Go without much ado, I
remain firmly a believer in the utility of functions and libraries in building solid systems, as not repeating the
same bit of code all over the code base lowers the probability of bugs.  Unless performance is a primary consideration,
I will likely use this package myself.

Also, if you are learning Go, you may find the changes between the original and final versions of interest.

-

bushwick, nyc
April 2012