package problems

// Deadlock: blocking of concurrent threads that have a mutual dependency that is impossible to be satisfied.

import (
	"testing"
)

// block on the wait channel before writing to the done and ctl channels
func deadlock(wait chan bool, done chan bool, ctl chan bool) {
	<-wait
	done <- true
	ctl <- true
}

func Test_Deadlock(t *testing.T) {
	c1 := make(chan bool)
	c2 := make(chan bool)
	ctl := make(chan bool)

	// start mutual dependent go routines
	go deadlock(c1, c2, ctl)
	go deadlock(c2, c1, ctl)

	<-ctl
	<-ctl
}
