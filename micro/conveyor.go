package micro

import (
	"github.com/micro/go-micro/broker"
	"runtime"
	"sync/atomic"

	"github.com/leolara/conveyor"
)

func NewBrokerFromMicroBroker(broker broker.Broker) conveyor.Broker {
	return &microConveyor{
		broker: broker,
	}
}

type microConveyor struct {
	broker broker.Broker
}

func (b *microConveyor) Subscribe(target string, options ...interface{}) <-chan conveyor.Subscription {
	sub := &subscription{ch: make(chan conveyor.ReceiveEnvelope)}
	sub.parent, sub.err = b.broker.Subscribe(target, sub.handler)

	ch := make(chan conveyor.Subscription, 1)
	ch <- sub
	close(ch)
	return ch
}

func (b *microConveyor) Publish(target string, msgs <-chan conveyor.SendEnvelop, options ...interface{}) {
	go func() {
		for msg := range msgs {
			bmsg := broker.Message{Body: msg.Body()}
			b.broker.Publish(target, &bmsg)
			go func () {msg.Error() <- nil}()
		}
	}()
}

type subscription struct {
	parent broker.Subscriber

	ch chan conveyor.ReceiveEnvelope

	err error
	closing uint32
}

func (s *subscription) Receive() <-chan conveyor.ReceiveEnvelope {
	return s.ch
}

func (s *subscription) Unsubscribe() {
	s.parent.Unsubscribe()

	atomic.StoreUint32(&s.closing, 1)
	runtime.Gosched()

	for {
		select {
		case <-s.ch:
			continue
		default:
			close(s.ch)
			return
		}
	}
}

func (s *subscription) Error() error {
	return s.err
}

func (s *subscription) handler(pub broker.Publication) error {
	if atomic.LoadUint32(&s.closing) != 0 {
		return nil
	}

	ack := make(chan interface{})
	msg := conveyor.NewReceiveEnvelopCopy(pub.Message().Body, ack)

	go func() {
		s.ch <- msg
		<-ack
		close(ack)
		pub.Ack()
	}()

	return nil
}

var _ conveyor.Subscription = (*subscription)(nil)
