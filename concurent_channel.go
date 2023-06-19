package asutils

import (
	"fmt"
	"sync"
)

// TODO testing , wait and range options to

// ChannelWithConcurrentSenders ... generic channel with concurrent producers
type ChannelWithConcurrentSenders[T any] interface {
	DetachSender() error
	ROChannel() <-chan T
	AttachSender_err() (ChannelWithConcurrentSenders[T], error)
	AttachSender() ChannelWithConcurrentSenders[T]
	// AllocateChannel() chan interface{}
	Send(T)
	Wait()
}

// ChannelWithConcurrentSendersReceiver ... RO recover type for ChannelWithConcurrentSenders for type safety
// type <-chan T <-chan interface{}

// concurrentChannel ... concurrently used and closeable string channel
type concurrentChannel[T any] struct {
	channel chan T
	closed  bool
	once    sync.Once
	wg      sync.WaitGroup
}

// NewChannelWithConcurrentSenders ... create a ChannelWithConcurrentSenders
func NewChannelWithConcurrentSenders[T any](length int) ChannelWithConcurrentSenders[T] {
	var ret ChannelWithConcurrentSenders[T]
	if length != 0 {
		ret = &concurrentChannel[T]{channel: make(chan T, length), closed: false}
	} else {
		ret = &concurrentChannel[T]{channel: make(chan T), closed: false}
	}
	return ret

}

// AttachSender increase reference counter and returns the interface itself , ideal for forward to an another goroutine
func (c *concurrentChannel[T]) AttachSender_err() (ChannelWithConcurrentSenders[T], error) {
	if !c.closed {
		c.wg.Add(1)
	} else {
		return nil, fmt.Errorf("channel already closed")
	}
	return c, nil
}

// AttachSender increase reference counter and returns the interface itself , ideal for forward to an another goroutine
func (c *concurrentChannel[T]) AttachSender() ChannelWithConcurrentSenders[T] {
	if !c.closed {
		c.wg.Add(1)
		return c
	} else {
		fmt.Println("ERROR:channel already closed")
		return nil
	}
	
}

// // AllocateChannel .. Allocate a sender and get the cannel
// func (c *concurrentChannel) AllocateChannel() chan interface{} {
// 	if !c.closed {
// 		c.wg.Add(1)
// 	} else {
// 		panic("channel already closed")
// 	}
// 	return c.channel
// }

// Wait ... wait until the all sender release the channel
func (c *concurrentChannel[T]) Wait() {
	c.wg.Wait()
}

// Send ... send an element to the channel
func (c *concurrentChannel[T]) Send(elem T) {
	c.channel <- elem
}

// ROChannel .. get receiver channel
func (c *concurrentChannel[T]) ROChannel() <-chan T {
	if c.closed {
		fmt.Println("ERROR:channel already closed")
		return nil
	}
	return c.channel
}

// DetachSender ...  Stop and close the channel
func (c *concurrentChannel[T]) DetachSender() error {
	if c.closed {
		return fmt.Errorf("channel already closed")
	}
	c.wg.Done()
	go func() {
		c.wg.Wait()
		c.once.Do(func() {
			close(c.channel)
			c.closed = true
		})
	}()
	return nil
}
