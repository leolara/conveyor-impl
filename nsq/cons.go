// Package nsq is an implementation for github.com/leolara/conveyor of system nsq.
// Currently it is a very thin wrapper around go-micro implementation, so tests are not included
package nsq

import (
	"github.com/leolara/conveyor-impl/internal"
	upstream "github.com/micro/go-plugins/broker/nsq"

	"github.com/leolara/conveyor"
)

// NewBroker creates a conveyor broker based on nsq
func NewBroker(addrs ...string) conveyor.Broker {
	return internal.NewFromMicroNew(upstream.NewBroker, addrs...)
}
