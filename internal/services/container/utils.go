package container

import (
	"log"
	"time"
)

// Timeout manages a timer with an associated action to be executed on expiration.
type Timeout struct {
	timer    *time.Timer
	duration time.Duration
	action   func()
}

// NewTimeout creates a new Timeout instance with the specified duration and action.
func NewTimeout(duration time.Duration, action func()) *Timeout {
	t := &Timeout{
		timer:    time.NewTimer(duration),
		duration: duration,
		action:   action,
	}
	return t
}

// StartTimer starts or resets the timer for the Timeout.
func (t *Timeout) StartTimer() {
	if !t.timer.Stop() {
		select {
		case <-t.timer.C:
			log.Fatalf("Timer should already been drained")
		default:
		}
	}
	t.timer.Reset(t.duration)

	go func() {
		<-t.timer.C
		if t.action != nil {
			t.action()
		}
	}()
}
