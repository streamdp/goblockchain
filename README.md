# goblockchain

<p align="center">
    <a href="https://github.com/streamdp/goblockchain/releases/latest">
        <img alt="GitHub release" src="https://img.shields.io/github/v/release/streamdp/goblockchain.svg?logo=github&style=flat-square">
    </a>
    <a href="https://goreportcard.com/report/github.com/streamdp/goblockchain">
        <img src="https://goreportcard.com/badge/github.com/streamdp/goblockchain" alt="Code Status" />
    </a>
</p>

This repo contains a simple/basic blockchain realisation in Go, with a basic code organization.
We use:
* gin-gonic/gin package to start and serve HTTP server
* crypto/sha1 to get SHA1 hashes

goblockchain use [Taskfile](https://dev.to/stack-labs/introduction-to-taskfile-a-makefile-alternative-h92) (a Makefile alternative). 

Please read the [Building Blockchain from Scratch in Python](https://python-scripts.com/blockchain) article in order to know more about this repository.

## Build the app
```bash
$ go build -o bin/goBlockChain internal/main.go
````
or
```bash
$ task build
````
## Run the app
```bash
$ ./bin/goBlockChain
```
or
```bash
$ task run
```
The default port is 8080, you can test the application in a browser or with curl:
```bash
$ curl 127.0.0.1:8080/chain
```
You can choose a different port and run more than one copy of goBlockChain on your local host.  For example:
```bash
$ ./bin/goBlockChain -port 8081
``` 
List of the endpoints:
* GET **/ping** _check node status_
* GET **/mine** _mine a new block_
* GET **/chain** _get the current state of the blockchain on a node_
* GET **/nodes/resolve** _get actual copy of the blockchain_ 
* GET **/mine/complexity/increase** _increase the difficulty of mining blocks_
* GET **/mine/complexity/decrease** _decrease the difficulty of mining blocks_
* POST **/transactions/new** _make a new transaction_
* POST **/nodes/register** _add a new node to the list of nodes_

Example of sending a POST request to add a new transaction to the blockchain:
```bash
$ curl -X POST -H "Content-Type: application/json" -d '{ "sender": "1914116639ac11ec83092c6fc90649b9", "recipient": "7e93670390396556d432206c1c3231fbb", "amount": 10}' "http://localhost:8080/transactions/new"
```