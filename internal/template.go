package internal

import (
	"github.com/leolara/conveyor"
	"github.com/leolara/conveyor-impl/micro"
	"github.com/micro/go-micro/broker"
)

type microBrokerConstructor func(opts ...broker.Option) broker.Broker

func NewFromMicroNew(cons microBrokerConstructor, addrs ...string) conveyor.Broker {
	return micro.NewBrokerFromMicroBroker(cons(broker.Addrs(addrs...)))
}
