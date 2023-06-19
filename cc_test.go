package asutils

import (
	"testing"
)

func TestConcurrentChannel(t *testing.T) {
	cc := NewChannelWithConcurrentSenders(10)
	cc.DetachSender()
}
