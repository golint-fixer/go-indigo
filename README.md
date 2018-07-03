# Indo Go

Core go implementation of Indo protocol.

[![Build Status](https://travis-ci.com/mitsukomegumi/indo-go.svg?branch=master)](https://travis-ci.com/mitsukomegumi/indo-go)

## Our Mission

Indo is the next generation of global decentralized information exchange--particularly that of monetary value--operating without central management. As a platform, Indo aims to make blockchain-based currencies more mainstream through Indo's absence of both transaction fees and excessive block wait times and intuitive integration APIs written entirely in Go. While the Indo platform aims to help publicize blockchain technologies, Indo also aims to help form and foster the inception and funding of individualism through the Indo token hosting platform, of which helps innovators and content-creators to get the funding they need while still remaining independent.

## True Decentralization

Indo's verification algorithms aren't based upon hashing power, but reliance and reputation, meaning that no one verification pool can gain dominance over the network.

## Contributing

Want to contribute to the development of Indo? Submit a pull request or contact one of our project admins!

## Building

Indo-go currently requires the latest build of the Go Language, which can be downloaded [here](https://golang.org/). After installing the Go Language, get the indo-go package by running

```bash
go get github.com/mitsukomegumi/indo-go
```

in your go-compatible command line of choice. The latest executable builds should be already be available, but feel free to build your own by running

```bash
go build
```

in the package's src folder.

### Running

Indo-go currently has 4 executable CLI flags used for testing and debugging:

```bash
go run main.go --relay
go run main.go --listen
go run main.go --fetch
go run main.go --host
```

### Dependencies

Indo-go requires golang's net package, of which can be acquired by running

```bash
go get golang.org/x/net
```

in any go-compatible terminal or command prompt.

## Code History

Historical data & initial source code can be found at the old Indo repositories: <https://github.com/IngeniousBlock/GeniCoin>, <https://github.com/mitsukomegumi/indocore>.
