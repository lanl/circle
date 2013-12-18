// This file demonstrates how to use the circle package.
package circle_test

import (
	"circle"
	"fmt"
)

var rank int // Our process's MPI rank

// createWork creates 10 units of "work" -- strings listing a rank and
// item number.
func createWork(q circle.Handle) {
	for i := 0; i < 10; i++ {
		work := fmt.Sprintf("RANK %d, ITEM %d", rank, i+1)
		if ok := q.Enqueue(work); !ok {
			panic("Enqueue")
		}
	}
}

// doWork processes one unit of "work" by dequeueing and outputting a string.
func doWork(q circle.Handle) {
	work, ok := q.Dequeue()
	if !ok {
		panic("Dequeue")
	}
	fmt.Printf("Rank %d is dequeueing %v\n", rank, work)
}

// This is an example of a complete program that uses the circle
// package.  It uses the low-level API (i.e., CallbackCreate,
// CallbackProcess, and Begin instead of ChannelBegin) and shows how
// to set various package options.
func Example() {
	// Initialize libcircle.
	rank = circle.Initialize()
	defer circle.Finalize()
	circle.EnableLogging(circle.LogErr) // Make libcircle a little quieter than normal.

	// Contrast the output when the following is uncommented (and
	// multiple MPI processes are used).
	//
	// circle.SetOptions(circle.CreateGlobal)

	// Create and execute some work.
	circle.CallbackCreate(createWork)
	circle.CallbackProcess(doWork)
	circle.Begin()
}
