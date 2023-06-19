package asutils

import (
	"fmt"
	"sync"
)

// TODO testing , wait and range options to

// ChannelWithConcurrentSenders ... generic channel with concurrent producers
type ChannelWithConcurrentSenders interface {
	DetachSender() error
	ROChannel() (ChannelWithConcurrentSendersReceiver, error)
	AttachSender() (ChannelWithConcurrentSenders, error)
	// AllocateChannel() chan interface{}
	Send(interface{})
	Wait()
}

// ChannelWithConcurrentSendersReceiver ... RO recover type for ChannelWithConcurrentSenders for type safety
type ChannelWithConcurrentSendersReceiver <-chan interface{}

// concurrentChannel ... concurrently used and closeable string channel
type concurrentChannel struct {
	channel chan interface{}
	closed  bool
	once    sync.Once
	wg      sync.WaitGroup
}

// NewChannelWithConcurrentSenders ... create a ChannelWithConcurrentSenders
func NewChannelWithConcurrentSenders(length int) ChannelWithConcurrentSenders {
	var ret ChannelWithConcurrentSenders
	if length != 0 {
		ret = &concurrentChannel{channel: make(chan interface{}, length), closed: false}
	} else {
		ret = &concurrentChannel{channel: make(chan interface{}), closed: false}
	}
	return ret

}

// AttachSender increase reference counter and returns the interface itself , ideal for forward to an another goroutine
func (c *concurrentChannel) AttachSender() (ChannelWithConcurrentSenders, error) {
	if !c.closed {
		c.wg.Add(1)
	} else {
		return nil, fmt.Errorf("channel already closed")
	}
	return c, nil
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
func (c *concurrentChannel) Wait() {
	c.wg.Wait()
}

// Send ... send an element to the channel
func (c *concurrentChannel) Send(elem interface{}) {
	c.channel <- elem
}

// ROChannel .. get receiver channel
func (c *concurrentChannel) ROChannel() (ChannelWithConcurrentSendersReceiver, error) {
	if c.closed {
		return nil, fmt.Errorf("channel already closed")
	}
	return c.channel, nil
}

// DetachSender ...  Stop and close the channel
func (c *concurrentChannel) DetachSender() error {
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
