[![Go Reference](https://pkg.go.dev/badge/github.com/sebundefined/thebus.svg)](https://pkg.go.dev/github.com/sebundefined/thebus)
[![Build Status](https://github.com/sebundefined/thebus/actions/workflows/release.yml/badge.svg)](https://github.com/sebudefined/thebus/actions/workflows/release.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/sebundefined/thebus)](https://goreportcard.com/report/github.com/sebundefined/thebus)

# thebus
**thebus** is a lightweight in-process pub/sub message bus for Go.
It provides a dead-simple API, fast fan-out, configurable delivery strategies (shared vs. cloned payloads), and sensible defaults for safety and performance.

## 🚀 Getting started

Just install **thebus** in your project by using the following command.

```shell
go get -u github.com/sebundefined/thebus
```

## 🔹 Quick Example

A minimal “hello world” pub/sub:

```go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/sebundefined/thebus"
)

func main() {
	// Create a new bus
	bus, _ := thebus.New()

	// Subscribe to a topic
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	sub, _ := bus.Subscribe(ctx, "hello")

	// Publish a message
	_, _ = bus.Publish("hello", []byte("world"))

	// Receive it
	msg := <-sub.Read()
	fmt.Printf("Got message on %s: %s\n", msg.Topic, msg.Payload)

	// Graceful shutdown
	_ = bus.Close()
}
```

## 🔧 With Custom Options

You can configure the bus globally:
```go
bus, _ := thebus.New(
	thebus.WithMaxTopics(100),
	thebus.WithCopyOnPublish(true), // safer if payloads are mutated after publish
)
```

Or customize subscribers

```go
sub, _ := bus.Subscribe(ctx, "foo",
	thebus.WithBufferSize(256),
	thebus.WithSendTimeout(100*time.Millisecond),
	thebus.WithStrategy(thebus.SubscriptionStrategyPayloadClonedPerSubscriber),
)
```

## ✨ Features

- ✅ Simple API (Publish, Subscribe, Unsubscribe, Stats, Close)
- 📦 Shared or cloned payload delivery strategies
- 🛡 Optional CopyOnPublish for safety against mutating payloads
- 📊 Backpressure & drop policies (DropIfFull, SendTimeout)
- 📉 Configurable limits (topics, subscribers per topic, buffer sizes)
- 🛑 Graceful shutdown with Close() and Unsubscribe()
- 🧪 Perfect for in-process events, simulations, and tests
- ⚡ Zero external deps (only stdlib crypto/rand)

## ❓Why this lib ?

We are using it at my current company for event driven within a single application. I decided to make it generic for publishing it to devs who need this
kind of in-process system (adding config, per subscriber strategy...)

## 🤔 Why choose thebus?
- Pure Go implementation — no broker required
- Lightweight alternative to Kafka/NATS when you just need local pub/sub
- Clean abstractions with good defaults
- Plays well with tests, mocks, and small services

## 🧪 Testing

Run the full suite:

```shell
go test -race ./...
```

With coverage:

```shell
go test -race -count=1 -coverprofile=coverage.out ./...
```

## Versioning

thebus follows Semantic Versioning.
Releases are published on GitHub and available via the Go module proxy.


## 🤝 Contributing

Contributions are welcome! Please check CONTRIBUTING.md.