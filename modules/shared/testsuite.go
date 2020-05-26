package shared

import (
	"runtime"
	"sync"
	"time"
)

var memStats runtime.MemStats

// Benchmark is the default structure to asses the performances
// of the functions
type Benchmark struct {
	Name     string
	N        int
	function func(b *Benchmark)

	mutex    sync.Mutex
	timerOn  bool
	failed   bool
	start    time.Time
	duration time.Duration
}

// StartTimer initializes the benchmark timer
func (b *Benchmark) StartTimer() {
	if !b.timerOn {
		b.start = time.Now()
		b.timerOn = true
	}
}

// StopTimer cancels the benchmark timer and calculates the duration
func (b *Benchmark) StopTimer() {
	if b.timerOn {
		b.duration += time.Now().Sub(b.start)
		b.timerOn = false
	}
}
