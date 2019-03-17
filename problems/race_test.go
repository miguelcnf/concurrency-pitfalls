package problems

// Race Conditions: concurrent threads changing shared data without sufficient protection mechanisms.

import (
	"testing"
)

const concurrency = 1000

// shared memory
var value = 0

// change the value of unprotected shared memory
func increment(done chan bool) {
	value++
	done <- true
}

func Test_RaceCondition(t *testing.T) {
	done := make(chan bool, concurrency)

	// start N go routines
	for i := 0; i < concurrency; i++ {
		go increment(done)
	}

	// wait for all go routines to finish
	for i := 0; i < concurrency; i++ {
		<-done
	}

	// assert that writes were lost due to race conditions
	if value != concurrency {
		t.Logf("expected value to be: %v but was: %v", concurrency, value)
		t.Fatal()
	}
}
