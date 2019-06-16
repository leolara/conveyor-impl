// Package mqtt is an implementation for github.com/leolara/conveyor of system mqtt.
// Currently it is a very thin wrapper around go-micro implementation, so tests are not included
package mqtt

import (
	"github.com/leolara/conveyor-impl/internal"
	upstream "github.com/micro/go-plugins/broker/mqtt"

	"github.com/leolara/conveyor"
)

// NewBroker creates a conveyor broker based on mqtt
func NewBroker(addrs ...string) conveyor.Broker {
	return internal.NewFromMicroNew(upstream.NewBroker, addrs...)
}
