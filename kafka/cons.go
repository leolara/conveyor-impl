// Package kafka is an implementation for github.com/leolara/conveyor of system kafka.
// Currently it is a very thin wrapper around go-micro implementation, so tests are not included
package kafka

import (
	"github.com/leolara/conveyor-impl/internal"
	upstream "github.com/micro/go-plugins/broker/kafka"

	"github.com/leolara/conveyor"
)

// NewBroker creates a conveyor broker based on kafka
func NewBroker(addrs ...string) conveyor.Broker {
	return internal.NewFromMicroNew(upstream.NewBroker, addrs...)
}
