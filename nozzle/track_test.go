package nozzle

import (
	"sync"
	"testing"
)

func TestResponseTimes(t *testing.T) {
	rt := &responseTimes{Mutex: &sync.RWMutex{}, Times: make([]int, 0), Index: 0, Ready: false, LastUpdate: int64(0)}
	if rt.Average() != 0 {
		t.Errorf("Expected average to be 0 but is %d", rt.Average())
	}

	rt.Times = append(rt.Times, 1)
	if rt.Average() != 1 {
		t.Errorf("Expected average to be 1 but is %d", rt.Average())
	}
	rt.Times = append(rt.Times, 2)
	if rt.Average() != 1.5 {
		t.Errorf("Expected average to be 1.5 but is %d", rt.Average())
	}
}

func TestTracker(t *testing.T) {
	tracker := NewTracker()
	url1 := "/foo"
	_, ok := tracker.lookupTimes(url1)
	if ok {
		t.Errorf("Expected not to find url %s in tracker")
	}

	tracker.AddResponse(url1, 1)
	rt, ok := tracker.lookupTimes(url1)
	if !ok {
		t.Fatalf("Expected the url with key %s to exist in tracker", url1)
	}
	if rt.Index != 1 {
		t.Errorf("Expected length of reponse times for %s to be 1 but was %d", url1, rt.Index)
	}

	_, ready := tracker.Average(url1)
	if ready {
		t.Fatalf("Should not be ready to get average for url %s", url1)
	}
	for x := 0; x < MinResponseTimes-1; x++ {
		tracker.AddResponse(url1, 1)
	}
	a, ready := tracker.Average(url1)
	if !ready {
		t.Fatalf("Expected tracker to be ready added enough data for url %s", url1)
	}
	if a != 1 {
		t.Errorf("Calculation is off average should be 1 but was %f", a)
	}
	tracker.AddResponse(url1, 100)
	a, ready = tracker.Average(url1)
	if !ready {
		t.Fatalf("Expected tracker to be ready added enough data for url %s", url1)
	}
	if a != 10 {
		t.Errorf("Calculation is off average should be 10 but was %f", a)
	}

	for x := 0; x < MaxResponseTimes; x++ {
		tracker.AddResponse(url1, 100)
	}
	a, ready = tracker.Average(url1)
	if !ready {
		t.Fatalf("Expected tracker to be ready added enough data for url %s", url1)
	}
	if a != 100 {
		t.Errorf("Calculation is off average should be 10 but was %f", a)
	}

}

func BenchmarkTracker(b *testing.B) {
	tracker := NewTracker()
	url1 := "/foo"
	for x := 0; x < MaxResponseTimes; x++ {
		tracker.AddResponse(url1, 100)
		tracker.Average(url1)
	}

}
