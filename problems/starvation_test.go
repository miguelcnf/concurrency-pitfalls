package problems

// Starvation: threads are unable to progress due to greedy or high priority threads blocking shared resources.

import (
	"sync"
	"testing"
	"time"
)

// shared resource
var l = &lock{}

type lock struct {
	m sync.Mutex
}

var thr = &threads{}

type threads struct {
	created  int
	finished int
	m        sync.Mutex
}

func (t *threads) create() {
	t.m.Lock()
	t.created++
	t.m.Unlock()
}

func (t *threads) finish() {
	t.m.Lock()
	t.finished++
	t.m.Unlock()
}

func (t *threads) values() (int, int) {
	t.m.Lock()
	defer t.m.Unlock()
	return t.created, t.finished
}

// greedy function that locks the shared resource for a long time
func greedy() {
	l.m.Lock()
	time.Sleep(2 * time.Second)
	l.m.Unlock()
}

// normal function that locks the shared resource for a small amount of time
func normal(finish chan bool) {
	l.m.Lock()
	time.Sleep(100 * time.Millisecond)
	finish <- true
	l.m.Unlock()
}

func Test_Starvation(t *testing.T) {
	create := make(chan bool)
	finish := make(chan bool)

	// register created/finished thread counters
	go func() {
		for {
			select {
			case <-create:
				thr.create()
			case <-finish:
				thr.finish()
			}
		}
	}()

	threads := time.Tick(100 * time.Millisecond)
	stop := time.After(5 * time.Second)
	greed := time.After(2 * time.Second)

loop:
	for {
		select {
		case <-greed:
			// start a single greedy go routine after 2s
			go greedy()
		case <-threads:
			// start a new normal go routine every 100ms
			go normal(finish)
			create <- true
		case <-stop:
			// wait 100ms to let normal go routines finish and break the loop after 5s
			time.Sleep(100 * time.Millisecond)
			break loop
		}
	}

	c, f := thr.values()
	if c != f {
		// assert that number of finished threads is different than the number of created threads
		t.Logf("number of created threads: %v", c)
		t.Logf("number of finished threads: %v", f)
		t.Fatal()
	}
}
