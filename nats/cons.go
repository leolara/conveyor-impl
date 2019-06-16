// Package nats is an implementation for github.com/leolara/conveyor of system nats.
// Currently it is a very thin wrapper around go-micro implementation, so tests are not included
package nats

import (
	"github.com/leolara/conveyor-impl/internal"
	upstream "github.com/micro/go-plugins/broker/nats"

	"github.com/leolara/conveyor"
)

// NewBroker creates a conveyor broker based on nats
func NewBroker(addrs ...string) conveyor.Broker {
	return internal.NewFromMicroNew(upstream.NewBroker, addrs...)
}
