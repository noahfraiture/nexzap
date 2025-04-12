package container

import (
	"log"
	"time"
)

// Timeout manages a timer with an associated action to be executed on expiration.
type Timeout struct {
	timer    *time.Timer
	duration time.Duration
	stop     chan any
	action   func()
}

// NewTimeout creates a new Timeout instance with the specified duration and action.
func NewTimeout(duration time.Duration, action func()) *Timeout {
	t := &Timeout{
		timer:    time.NewTimer(duration),
		duration: duration,
		stop:     make(chan any),
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
		log.Printf("Starting timer for duration %v", t.duration)
		selectStart := time.Now()
		select {
		case <-t.timer.C:
			log.Printf("Timer expired after %v, executing action", time.Since(selectStart))
			if t.action != nil {
				t.action()
			}
		case <-t.stop:
			log.Printf("Timer stopped after %v", time.Since(selectStart))
		}
	}()
}
