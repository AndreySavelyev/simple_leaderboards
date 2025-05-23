# events mock service

This if fairly straightforward service.
It consists of one file and a few functions.

Main aim of this service is to generate a consistent flow of event emulating some betting activity with occasional bursts.

To do that, there exists an algorithm named TokenBucket algotirm.
It's implemented in the https://pkg.go.dev/golang.org/x/time/rate library

More on this here https://en.wikipedia.org/wiki/Token_bucket

Generally, we have a 'bucket' of a fixed size(though in our case there is no overflow) which get's filled on every tick(by the library).
We're trying to use X number of tokens, if there is such number of tokens in the bucket, we 'spend' them and generate necessary number of events.
If there isn't enough tokens, we wait for a time given to us by the library.


This service is stateless and generates events without storing them anywhere. We assume that the transport library is reliable and highly available(imagine kafka instead of redis here).
Communication is done via Redis Pub/Sub mechanism.

Some events' attributes such are Currencies, Game names, Distributors, Studios are taken out of the hardcoded list for simplicity.
There is a duplication of information between this service and the consumer regarging the currencies and probably something more centralied could've been preferred, but for the ease of implementation, I decided to go this way.

## How to run
Open a terminal tab
```
> go build main.go && ./main -users 50
```
The service will that logging some info about the events generated.
