/*
This file provides a higher-level, more Go-like API to the libcircle library.
*/

package circle

// ChannelBegin replaces CallbackCreate, CallbackProcess, and Begin
// with a channel-based interface.  The caller is expected to write
// work into putWork, close the channel, then read work from getWork.
// An implication is that no new work can be created after the initial
// set of work is written into putWork.  Use the lower-level API
// (CallbackCreate, CallbackProcess, and Begin) if workers need to be
// able to enqueue new work.
func ChannelBegin() (putWork chan<- string, getWork <-chan string) {
	// Create our input and output channels.
	toCircle := make(chan string)
	fromCircle := make(chan string)

	// Queue processing runs in the background.
	go func() {
		// Register a creation callback that reads from
		// toCircle and writes to the queue.
		CallbackCreate(func(q Handle) {
			for work := range toCircle {
				q.Enqueue(work)
			}
		})

		// Register a processing callback that reads from the
		// queue and writes to fromCircle.
		CallbackProcess(func(q Handle) {
			work, _ := q.Dequeue()
			fromCircle <- work
		})

		// Begin processing the queue.
		Begin()
		close(fromCircle)
	}()

	// Return our channels to the user.
	putWork = toCircle
	getWork = fromCircle
	return
}
