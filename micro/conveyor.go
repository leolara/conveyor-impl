package micro

import (
	"sync"

	"github.com/micro/go-micro/broker"

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
	sub := &subscription{ch: make(chan conveyor.ReceiveEnvelope), stopCh: make(chan interface{})}
	sub.parent, sub.err = b.broker.Subscribe(target, sub.handler)

	ch := make(chan conveyor.Subscription, 1)
	ch <- sub
	close(ch)
	return ch
}

func (b *microConveyor) Publish(target string, msgs <-chan conveyor.SendEnvelop, options ...interface{}) {
	go func() {
		for msg := range msgs {
			go func(msg conveyor.SendEnvelop) {
				bmsg := broker.Message{Body: msg.Body()}
				b.broker.Publish(target, &bmsg)
				msg.Error() <- nil
			}(msg)
		}
	}()
}

type subscription struct {
	parent broker.Subscriber

	ch chan conveyor.ReceiveEnvelope

	err            error
	stopCh         chan interface{}
	writersWG      sync.WaitGroup
	writersWGMutex sync.Mutex
}

func (s *subscription) Receive() <-chan conveyor.ReceiveEnvelope {
	return s.ch
}

func (s *subscription) Unsubscribe() {
	s.parent.Unsubscribe()

	close(s.stopCh)

	s.writersWGMutex.Lock()
	s.writersWG.Wait()
	s.writersWGMutex.Unlock()

	close(s.ch)
}

func (s *subscription) Error() error {
	return s.err
}

func (s *subscription) handler(pub broker.Publication) error {
	go func() {
		s.writersWGMutex.Lock()
		s.writersWG.Add(1)
		s.writersWGMutex.Unlock()
		defer s.writersWG.Done()

		select {
		case <-s.stopCh:
			return
		default:
		}

		ack := make(chan interface{})
		msg := conveyor.NewReceiveEnvelopCopy(pub.Message().Body, ack)

		select {
		case <-s.stopCh:
		case s.ch <- msg:
			go func(ack chan interface{}) {
				// pass through ack
				<-ack
				close(ack)
				pub.Ack()
			}(ack)
		}
	}()

	return nil
}

var _ conveyor.Subscription = (*subscription)(nil)
