// Package sqs is an implementation for github.com/leolara/conveyor of system sqs.
// Currently it is a very thin wrapper around go-micro implementation, so tests are not included
package sqs

import (
	"github.com/leolara/conveyor-impl/internal"
	upstream "github.com/micro/go-plugins/broker/sqs"

	"github.com/leolara/conveyor"
)

// NewBroker creates a conveyor broker based on sqs
func NewBroker(addrs ...string) conveyor.Broker {
	return internal.NewFromMicroNew(upstream.NewBroker, addrs...)
}
