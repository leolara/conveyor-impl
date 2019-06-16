// Package rabbitmq is an implementation for github.com/leolara/conveyor of system rabbitmq.
// Currently it is a very thin wrapper around go-micro implementation, so tests are not included
package rabbitmq

import (
	"github.com/leolara/conveyor-impl/internal"
	upstream "github.com/micro/go-plugins/broker/rabbitmq"

	"github.com/leolara/conveyor"
)

// NewBroker creates a conveyor broker based on rabbitmq
func NewBroker(addrs ...string) conveyor.Broker {
	return internal.NewFromMicroNew(upstream.NewBroker, addrs...)
}
