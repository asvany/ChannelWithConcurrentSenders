package asutils

import (
	"sync"
)

// TODO testing , wait and range options to

// ChannelWithConcurrentSenders ... generic channel with concurrent producers
type ChannelWithConcurrentSenders interface {
	Close()
	ROChannel() ChannelWithConcurrentSendersReceiver
	Forward() ChannelWithConcurrentSenders
	AllocateChannel() chan interface{}
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
	if length != 0 {
		return &concurrentChannel{channel: make(chan interface{}, length), closed: false}
	}
	return &concurrentChannel{channel: make(chan interface{}), closed: false}
}

// Forward increase reference counter and returns the interface itself , ideal for forward to an another goroutine
func (c *concurrentChannel) Forward() ChannelWithConcurrentSenders {
	if !c.closed {
		c.wg.Add(1)
	} else {
		panic("channel already closed")
	}
	return c
}

// Wait ... wait until the all sender release the channel
func (c *concurrentChannel) Wait() {
	c.wg.Wait()
}

// Send ... send an element to the channel
func (c *concurrentChannel) Send(elem interface{}) {
	c.channel <- elem
}

// AllocateChannel .. Allocate a sender and get the cannel
func (c *concurrentChannel) AllocateChannel() chan interface{} {
	if !c.closed {
		c.wg.Add(1)
	} else {
		panic("channel already closed")
	}
	return c.channel
}

// ROChannel .. get receiver channel
func (c *concurrentChannel) ROChannel() ChannelWithConcurrentSendersReceiver {
	if c.closed {
		panic("channel already closed")
	}
	return c.channel
}

// Close ...  Stop and close the channel
func (c *concurrentChannel) Close() {
	c.wg.Done()
	go func() {
		c.wg.Wait()
		c.once.Do(func() {
			close(c.channel)
			c.closed = true
		})
	}()
}
