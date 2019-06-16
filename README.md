# Broker implementations for github.com/leolara/conveyor

[![Build Status](https://travis-ci.org/leolara/conveyor-impl.svg?branch=master)](https://travis-ci.org/leolara/conveyor-impl)
[![docs](https://godoc.org/github.com/leolara/conveyor-impl?status.svg)](http://godoc.org/github.com/leolara/conveyor-impl)

At this repo you can find implementations for different messaging systems that are compatible with Conveyor, so you 
can use them in a Go idiomatic and asynchronous way using Go channels.

Implementations:
 + Kafka
 + MQTT
 + NATS
 + NSQ
 + RabbitMQ
 + Redis
 + SQS
 + STAN (NATS streaming)
 + STOMP
 + Micro, a wrapper around go-micro project brokers that convert them 
 
 Currently, there are all implemented using a small wrapper around Go micro broker libraries.
 
 For more information about conveyor: [github.com/leolara/conveyor](https://github.com/leolara/conveyor)
 