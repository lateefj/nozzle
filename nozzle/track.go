// This module keeps track of average response times over a specific period of time
package nozzle

import (
	"sync"
	"time"
)

const (
	MinResponseTimes = 10
	MaxResponseTimes = 100
)

type responseTimes struct {
	Mutex      *sync.RWMutex
	Times      []int
	Index      int
	Ready      bool
	LastUpdate int64
}

func (rt *responseTimes) Average() float64 {
	total := 0
	for _, v := range rt.Times {
		total += v
	}
	return float64(total) / float64(len(rt.Times))
}

type Tracker struct {
	// Keep locks concurrent
	mutex *sync.RWMutex
	times map[string]*responseTimes
}

func NewTracker() *Tracker {
	return &Tracker{&sync.RWMutex{}, make(map[string]*responseTimes)}
}

// Provides read locks around the map
func (t *Tracker) lookupTimes(r string) (*responseTimes, bool) {
	t.mutex.RLock()
	defer t.mutex.Unlock()
	rt, ok := t.times[r]
	return rt, ok
}

// Gets the average response times for a key
func (t *Tracker) Average(r string) (bool, float64) {
	ready := false
	average := float64(0)
	rTimes, ok := t.lookupTimes(r)
	if ok {
		rTimes.Mutex.RLock()
		ready = rTimes.Ready
		average = rTimes.Average()
		rTimes.Mutex.Unlock()
	}
	return ready, average
}

// Add a response time to a key
func (t *Tracker) AddResponse(r string, rt int) {
	rTimes, ok := t.lookupTimes(r)
	// If can't find the reponse times then create a new one for this
	if !ok {
		rTimes = &responseTimes{Times: make([]int, MaxResponseTimes), Index: 0, Ready: false, LastUpdate: int64(0)}
		// Set a response time for a specific key
		t.mutex.Lock()
		t.times[r] = rTimes
		t.mutex.Unlock()
	}
	// Grab reponse times lock to do some writing with
	rTimes.Mutex.Lock()
	defer rTimes.Mutex.Unlock()
	// Store the time
	rTimes.Times[rTimes.Index] = rt

	// Track the last time the response time was updated
	rTimes.LastUpdate = time.Now().UnixNano()
	// Set the response time ready if at min number
	if len(rTimes.Times) > MinResponseTimes {
		rTimes.Ready = true
	}
	// Increment the index
	rTimes.Index += 1
	// Loop back if the index is at the max
	if rTimes.Index >= MaxResponseTimes {
		rTimes.Index = 0
	}
}
