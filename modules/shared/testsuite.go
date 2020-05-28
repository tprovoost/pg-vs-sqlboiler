package shared

import (
	"fmt"
	"time"
)

var DebugMode = false

// Benchmark is the default structure to asses the performances
// of the functions
type Benchmark struct {
	N int

	timerOn  bool
	Failed   bool
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

// GetDuration is an accessor to duration property
func (b *Benchmark) GetDuration() time.Duration {
	return b.duration
}

// BenchmarkSuite is a structure holding all benchmarks
type BenchmarkSuite struct {
	Insert       []Benchmark
	ReadOne      []Benchmark
	ReadAll      []Benchmark
	FetchIn      []Benchmark
	ComplexQuery []Benchmark
}

// Print displays all results in an understandable manner
func (suite *BenchmarkSuite) Print() {
	fmt.Println("Function\tSQL Boiler\tPG")
	fmt.Println("----------------------------------------")
	fmt.Print("Insert:\t\t")
	printBenchmarks(suite.Insert)

	fmt.Print("\nReadOne:\t")
	printBenchmarks(suite.ReadOne)

	fmt.Print("\nFetch in:\t")
	printBenchmarks(suite.FetchIn)

	fmt.Println()
}

func printBenchmarks(bs []Benchmark) {
	for _, b := range bs {
		fmt.Printf("%d\t\t", b.duration.Microseconds())
	}
}
