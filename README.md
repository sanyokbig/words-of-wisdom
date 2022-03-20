# Words of Wisdom
## Task description
Test task for Server Engineer

Design and implement “Word of Wisdom” tcp server.              
• TCP server should be protected from DDOS attacks with the Prof of Work (https://en.wikipedia.org/wiki/Proof_of_work), the challenge-response protocol should be used.              
• The choice of the POW algorithm should be explained.              
• After Prof Of Work verification, server should send one of the quotes from “word of wisdom” book or any other collection of the quotes.              
• Docker file should be provided both for the server and for the client that solves the POW challenge.

## How to use
### Build
```shell
$ make build-server
$ make build-client
```

### Run
```shell
$ make run-server
$ make run-client
```

### Clean Up
```shell
$ make clean-up
```

## Implementation
Implemented PoW design is based on a [Moderately Hard, Memory-bound Functions](https://users.soe.ucsc.edu/~abadi/Papers/memory-longer-acm.pdf) work.

Memory-bound PoW is chosen as memory access performance is less sensitive to hardware and should work fine on both low and high-end hardware.

General design matches a proposed one where Server generates a random _x0_ and applies _F()_ to it _k_ times resulting in _xk_.

Client knows all information about the parameters of algorithm except the _x0_ and is expected to try all different paths towards _x0_.

When _x0_ is found, Client compares a checksum of a sequence to the checksum of a valid sequence received from the Server. When checksum matches, solution is found. Is checksums don't match, client goes for another sequence until valid is found.

### Function F()
Implementation of _F()_ will greatly affect difficulty and efficiency of an algorithm.
It is desirable that there are _x_ and _x'_, where _F(x)=F(x')_. This requires client to traverse both paths to check sequences, increasing required work.

In addition, it's required that calling inverted F() is slower that accessing inversion table, encouraging client to use memory instead of CPU. 

Current implementation of F() is pretty basic and probably could be reversed faster than accessing a memory. But it can be swapped for a different, more difficult function if needed.  

### Challenge difficulty
Difficulty of this PoW can be configured with two parameters _k_ and _n_.

_k_ represents a number of times _F()_ is applied. Increasing this parameter will result in longer sequences and will affect both Client and Server

_n_ represents a range of possible values, which will be in a range [0, 2^n). Increasing it will mostly affect Client and not a Server, since Client will have to generate a larger inversion table and process more sequences, while Server will just generate a larger numbers. 


### Performance
With `k = 64` and  `n = 21`, there are following results.

On a first dev machine:
- Server prepared challenge in 55µs
- Client solved challenge in 568ms

On a second, significantly less powerful dev machine:
- Server prepared challenge in 72µs
- Client solved challenge in 660ms
