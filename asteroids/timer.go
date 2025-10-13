// File timer.go defines a lightweight, tick-based timer abstraction.
// Timers are central to spawning, pacing, cooldowns, and animations
// throughout the game, synchronized to Ebiten’s tick rate rather than wall time.
package asteroids

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

// Timer represents a simple counter-based timer that progresses
// in sync with Ebiten’s ticks (frames). Once currentTicks >= targetTicks,
// the timer is considered "ready".
type Timer struct {
	currentTicks int // Elapsed tick count since last reset.
	targetTicks  int // Total ticks required before IsReady() is true.
}

// NewTimer returns a Timer for the specified duration.
//
// The provided duration (time.Duration) is converted to Ebiten ticks
// based on the target ticks-per-second (ebiten.TPS()). This ensures
// consistent timing behavior across platforms and frame rates.
func NewTimer(d time.Duration) *Timer {
	return &Timer{
		currentTicks: 0,
		targetTicks:  int(d.Milliseconds()) * ebiten.TPS() / 1000,
	}
}

// Update advances the timer by one tick, up to the target threshold.
//
// This should be called once per frame inside a scene or object’s Update().
func (t *Timer) Update() {
	if t.currentTicks < t.targetTicks {
		t.currentTicks++
	}
}

// IsReady returns true when the timer has reached or exceeded its target.
//
// Example:
//
//	if timer.IsReady() {
//	    spawnAlien()
//	    timer.Reset()
//	}
func (t *Timer) IsReady() bool {
	return t.currentTicks >= t.targetTicks
}

// Reset sets the timer’s tick count back to zero, restarting the interval.
func (t *Timer) Reset() {
	t.currentTicks = 0
}
