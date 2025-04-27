package container

import (
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
	t.StartTimer()
	return t
}

// StartTimer starts or resets the timer for the Timeout.
func (t *Timeout) StartTimer() {
	// Stop the timer if it's already running
	if !t.timer.Stop() {
		<-t.timer.C
	}
	// Reset the timer to start fresh
	t.timer.Reset(t.duration)

	go func() {
		<-t.timer.C
		if t.action != nil {
			t.action()
		}
	}()
}
