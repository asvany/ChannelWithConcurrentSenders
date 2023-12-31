package cc

import (
	"math/rand"
	"testing"
)

const noSender = 10
const noReceiver = 10
const noMessage = 20

func randomIntSender(cc ChannelWithConcurrentSenders[int]) {
	defer cc.DetachSender()
	for i := 0; i < noMessage; i++ {
		cc.Send(rand.Intn(100))
	}
}

func TestConcurrentChannel_Int_SingleReceiver(t *testing.T) {
	_cc := NewChannelWithConcurrentSenders[int](10)
	for i := 0; i < noSender; i++ {
		go randomIntSender(_cc.AttachSender())
	}
	counter := 0
	for i := range _cc.ROChannel() {
		counter++
		t.Log(i)
	}
	_cc.Wait()
	if counter != noSender*noMessage {
		t.Errorf("bad result: got %v, want %v", counter, noSender*noMessage)
	}
	
}

func TestConcurrentChanel_Int_MultipleReceivers(t *testing.T) {
	_cc := NewChannelWithConcurrentSenders[int](10)
	for i := 0; i < noSender; i++ {
		go randomIntSender(_cc.AttachSender())
	}
	counter := 0
	for i := 0; i < noReceiver; i++ {
		go func() {
			for i := range _cc.ROChannel() {
				counter++
				t.Log(i, counter)
			}
		}()
	}
	_cc.Wait()

	if counter != noSender*noMessage {
		t.Errorf("bad result: got %v, want %v", counter, noSender*noMessage)
	}

	
}
