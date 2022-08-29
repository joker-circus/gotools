package timer

import (
	"testing"
	"time"
)

// test initTimer function
func TestTimerInitTimer(t *testing.T) {
	// test nil Timer
	var nilTimer *time.Timer
	resNilTime := initTimer(nilTimer, 2*time.Second)
	if resNilTime == nil {
		t.Fatalf("Unexpected a nil. Expecting a Timer.")
	}

	// test the panic
	panicTimer := time.NewTimer(1 * time.Second)
	resPanicTimer := wrapInitTimer(panicTimer, 2*time.Second)
	if resPanicTimer != -1 {
		t.Fatalf("Expecting a panic for Timer, but nil")
	}
	// sleep enough time to test next timer
	time.Sleep(3 * time.Second)
}

func wrapInitTimer(t *time.Timer, timeout time.Duration) (ret int) {
	defer func() {
		if err := recover(); err != nil {
			ret = -1
		}
	}()
	res := initTimer(t, timeout)
	if res != nil {
		ret = 1
	}
	return ret
}

func TestTimerStopTimer(t *testing.T) {
	normalTimer := time.NewTimer(3 * time.Second)
	stopTimer(normalTimer)
	if normalTimer.Stop() {
		t.Fatalf("Expecting timer stopped, but it doesn't")
	}
}

func TestTimerAcquireTimer(t *testing.T) {
	normalTimer := AcquireTimer(2 * time.Second)
	if normalTimer == nil {
		t.Fatalf("Unexpected nil, expecting a timer")
	}
	ReleaseTimer(normalTimer)
}

func TestTimerReleaseTimer(t *testing.T) {
	normalTimer := AcquireTimer(2 * time.Second)
	ReleaseTimer(normalTimer)
	if normalTimer.Stop() {
		t.Fatalf("Expecting the timer is released.")
	}
}
