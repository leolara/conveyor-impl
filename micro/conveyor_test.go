package micro

import (
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/micro/go-micro/broker/memory"

	"github.com/leolara/conveyor"
)

func TestMicroMemoryBroker(t *testing.T) {
	mb := memory.NewBroker()
	mb.Connect()
	b := NewBrokerFromMicroBroker(mb)

	sub := <-b.Subscribe("testTopic")

	if sub.Error() != nil {
		t.Fatal(sub.Error())
	}

	pubChan := make(chan conveyor.SendEnvelop)
	pubChanErr := make(chan error)
	b.Publish("testTopic", pubChan)

	pubChan <- conveyor.NewSendEnvelop([]byte{24}, pubChanErr)

	select {
	case err := <-pubChanErr:
		if err != nil {
			t.Errorf("got publication error: %s", err)
		}
	case <-time.After(10 * time.Millisecond):
		t.Error("Did not receive empty error")
	}

	select {
	case envelope, ok := <-sub.Receive():
		if !ok {
			t.Error("channel closed")
		}
		if len(envelope.Body()) != 1 && envelope.Body()[0] != 24 {
			t.Error("received wrong data")
		}
		envelope.Ack() <- nil
	case <-time.After(10 * time.Millisecond):
		t.Error("Did not receive message")
	}

	pubChan <- conveyor.NewSendEnvelop([]byte{42}, pubChanErr)
	<-pubChanErr
	pubChan <- conveyor.NewSendEnvelop([]byte{42}, pubChanErr)
	<-pubChanErr

	sub.Unsubscribe()
	pubChan <- conveyor.NewSendEnvelop([]byte{42}, pubChanErr)
	<-pubChanErr
	runtime.Gosched()

	select {
	case _, ok :=<-sub.Receive():
		if ok {
			t.Error("shouldn't receive anything")
		}
	case <-time.After(10 * time.Millisecond):
		// OK
	}

	pubChan <- conveyor.NewSendEnvelop([]byte{42}, pubChanErr)
	runtime.Gosched()

	select {
	case _, ok :=<-sub.Receive():
		if ok {
			t.Error("shouldn't receive anything")
		}
	case <-time.After(10 * time.Millisecond):
		// OK
	}
}

func TestMicroMemoryBrokerRace(t *testing.T) {
	producer := func(b conveyor.Broker, n int, wg *sync.WaitGroup) {
		defer wg.Done()
		wg.Add(1)

		pubChan := make(chan conveyor.SendEnvelop)
		pubChanErr := make(chan error)
		b.Publish("testTopic", pubChan)

		for i:=1; i<n; i++ {
			pubChan <- conveyor.NewSendEnvelop([]byte{24}, pubChanErr)
			err := <-pubChanErr
			if err != nil {
				t.Error(err)
			}
			runtime.Gosched()
		}

		runtime.Gosched()
		close(pubChan)
		runtime.Gosched()
	}

	consumer := func(b conveyor.Broker, n int, wg *sync.WaitGroup) {
		defer wg.Done()
		wg.Add(1)

		sub := <-b.Subscribe("testTopic")

		if sub.Error() != nil {
			t.Fatal(sub.Error())
		}

		for i:=1; i<n; i++ {
			select {
			case envelope, ok := <-sub.Receive():
				if !ok {
					t.Error("channel closed")
				}
				if len(envelope.Body()) != 1 && envelope.Body()[0] != 24 {
					t.Error("received wrong data")
				}
				envelope.Ack() <- nil
			case <-time.After(time.Second):
				t.Error("Did not receive message")
			}
		}

		sub.Unsubscribe()

		for i:=1; i<n; i++ {
			select {
			case _, ok := <-sub.Receive():
				if ok {
					t.Error("shouldn't receive anything")
				}
			case <-time.After(10 * time.Millisecond):
				// OK
			}
		}
	}

	for i:=1; i<100; i++ {
		mb := memory.NewBroker()
		mb.Connect()
		b := NewBrokerFromMicroBroker(mb)

		var wg sync.WaitGroup
		go consumer(b, 10, &wg)
		go consumer(b, 20, &wg)

		go producer(b, 10, &wg)
		go producer(b, 30, &wg)
		go producer(b, 30, &wg)
		wg.Wait()
	}
}
