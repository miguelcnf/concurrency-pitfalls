package problems

import (
	"runtime"
	"testing"
)

// Go Routines Leak: go routines are not garbage collected even when blocked on unreachable channels.

func Test_Leak(t *testing.T) {
	block := make(chan bool)

	// start N go routines
	for i := 0; i < 1000; i++ {
		// block on channel receive
		go func(c chan bool) {
			// channel is unreachable as there are no possible receivers
			c <- true
		}(block)
	}

	// number of live go routines
	liveGoRoutines := runtime.NumGoroutine()

	t.Logf("number of go routines before GC: %v", liveGoRoutines)

	// force GC
	runtime.GC()

	// assert that blocked go routines are not garbage collected when blocked on unreachable channels
	if liveGoRoutines == runtime.NumGoroutine() {
		t.Logf("expected number of go routines to have decreased after GC but was still: %v", runtime.NumGoroutine())

	}

}
