// Package redis is an implementation for github.com/leolara/conveyor of system redis.
// Currently it is a very thin wrapper around go-micro implementation, so tests are not included
package redis

import (
	"github.com/leolara/conveyor-impl/internal"
	upstream "github.com/micro/go-plugins/broker/redis"

	"github.com/leolara/conveyor"
)

// NewBroker creates a conveyor broker based on redis
func NewBroker(addrs ...string) conveyor.Broker {
	return internal.NewFromMicroNew(upstream.NewBroker, addrs...)
}
