# DDOS protection with Proof of Work


## Problem
Design and implement “Word of Wisdom” tcp server.
- TCP server should be protected from DDOS attacks with the Prof of Work (https://en.wikipedia.org/wiki/Proof_of_work), the challenge-response protocol should be used.
- The choice of the POW algorithm should be explained.
- After Prof Of Work verification, server should send one of the quotes from “word of wisdom” book or any other collection of the quotes.
- Docker file should be provided both for the server and for the client that solves the POW challenge

## Solution
- The PoW algorithm was chosen to be the same one used in bitcoin (hashcash) because it is simple and has proven to be stable over the years.
- Message types are json, for simplicity and compatibility with any programming language.

## How to run

### Server
```bash
make docker-build-server
make docker-run-server
```

### Client
```bash
make docker-build-client
make docker-run-client
```

## How to stop
```bash
make docker-down
```
