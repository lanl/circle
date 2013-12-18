// This file demonstrates how to use the higher-level, channel,
// interface to the circle package.
package circle_test

import (
	"encoding/json"
	"fmt"
	"github.com/losalamos/circle"
	"log"
)

// A Point contains x and y coordinates.
type Point struct {
	X, Y float64
}

// String implements the fmt.Stringer interface for pretty-printed output.
func (pt Point) String() string {
	return fmt.Sprintf("(%5.2f, %5.2f)", pt.X, pt.Y)
}

// Demonstrate how to use ChannelBegin to enqueue a bunch of Point
// objects then have remote workers dequeue and "process" (in this
// case, print) them.
func ExampleChannelBegin() {
	// Initialize libcircle.
	rank := circle.Initialize()
	defer circle.Finalize()

	// Create a pair of channels for writing work into the queue
	// and reading work from the queue.
	toQ, fromQ := circle.ChannelBegin()

	// Process 0 writes a bunch of Points into the queue.
	if rank == 0 {
		for j := 0; j < 5; j++ {
			for i := 0; i < 5; i++ {
				pt := Point{X: float64(i) * 1.23, Y: float64(j) * 4.56}
				enc, err := json.Marshal(pt)
				if err != nil {
					log.Fatalln(err)
				}
				toQ <- string(enc)
			}
		}
		close(toQ)
	}

	// All processes read Points from the queue and output them.
	for work := range fromQ {
		var pt Point
		if err := json.Unmarshal([]byte(work), &pt); err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("Dequeueing %s\n", pt)
	}
}
