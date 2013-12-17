/*
Package circle is a Go interface to the libcircle distributed-queue API.
libcircle is available from https://github.com/hpc/libcircle.
*/
package circle

/*
#cgo LDFLAGS: -lcircle
#include <libcircle.h>
#include <stdlib.h>
#include <stdio.h>

// Invoke libcircle's CIRCLE_cb_create() with a fixed callback function
// that will itself invoke the user-provided callback function.
void invoke_cb_create(void)
{
  extern void callUserCBCreate(CIRCLE_handle *);
  CIRCLE_cb_create((CIRCLE_cb)callUserCBCreate);
}

// Invoke libcircle's CIRCLE_cb_process() with a fixed callback function
// that will itself invoke the user-provided callback function.
void invoke_cb_process(void)
{
  extern void callUserCBProcess(CIRCLE_handle *);
  CIRCLE_cb_process((CIRCLE_cb)callUserCBProcess);
}

// Invoke an enqueue or dequeue C function.  (I wasn't able to get
// this to work without a helper function.)
int8_t invoke_queue_func(int8_t (*qfunc)(char *), char *elt)
{
  return qfunc(elt);
}

*/
import "C"

import (
	"os"
	"unsafe"
)

// We don't know a priori how much space to allocate to store a
// dequeued work item.  Set MaxWorkItemLength larger than the largest
// item you expect to dequeue.
var MaxWorkItemLength int = 1 << 20 // This size "ought to be enough for anybody".

// A Handle provides an interface to enqueue and dequeue libcircle work items.
type Handle interface {
	Enqueue(string) bool     // Enqueue a user-defined work item.  Return a success code.
	Dequeue() (string, bool) // Dequeue and return a user-defined work item plus a success code.
	LocalQueueSize() uint32  // Number of entries currently in the local queue
}

// A queue is a Go wrapper for a libcircle handle.  Its purpose is to
// let us define the Handle interface.
type queue struct {
	handle *C.CIRCLE_handle // libcircle internal handle type
}

// Enqueue a work item (a user-defined string) onto a queue.  Return a
// success code.  Note that the work item will be truncated at the
// first null character so binary strings are best avoided.
func (q queue) Enqueue(work string) (ok bool) {
	cstr := C.CString(work)
	defer C.free(unsafe.Pointer(cstr))
	ok = C.invoke_queue_func(q.handle.enqueue, cstr) >= 0
	return
}

// Dequeue a work item (a user-defined string) from a queue.  Return
// it plus a success code.
func (q queue) Dequeue() (str string, ok bool) {
	// This is obnoxious -- we don't know how large of a string to
	// allocate.  Allocate something fairly large and hope for the
	// best.
	buffer := make([]byte, MaxWorkItemLength)
	cstr := C.CString(string(buffer))
	defer C.free(unsafe.Pointer(cstr))

	// Dequeue the work item into the string, and return it.
	if ok = C.invoke_queue_func(q.handle.dequeue, cstr) >= 0; ok {
		str = C.GoString(cstr)
	}
	return
}

// LocalQueueSize returns the number of entries currently in the local
// work queue.
func (q queue) LocalQueueSize() (sz uint32) {
	type LQSFunc *func() C.uint32_t
	lqsFunc := LQSFunc(unsafe.Pointer(q.handle.local_queue_size))
	return uint32((*lqsFunc)())
}

// Initialize initializes libcircle and returns the current MPI rank.
func Initialize() (rank int) {
	argc := len(os.Args)
	argv := make([]*C.char, argc)
	for i, arg := range os.Args {
		argv[i] = C.CString(arg)
		defer C.free(unsafe.Pointer(argv[i]))
	}
	rank = int(C.CIRCLE_init(C.int(argc), &argv[0], C.int(DefaultFlags)))
	return
}

// A Flag is passed to SetOptions and controls libcircle's global behavior.
type Flag int32

// These constants can be ORed together to produce a Flag.
const (
	SplitRandom  = Flag(1 << iota) // Split work randomly.
	SplitEqual                     // Split work evenly.
	CreateGlobal                   // Call the creation callback on all processes.
	DefaultFlags = SplitEqual      // Default behavior is random work stealing.
)

// SetOptions sets libcircle's global behavior according to the
// inclusive-or of a set of flags.
func SetOptions(options Flag) {
	C.CIRCLE_set_options(C.int(options))
}

// A Callback is a user-provided function that libcircle will invoke
// as necessary.
type Callback func(Handle)

// Replicate parts of libcircle's internal CIRCLE_input_st structure
// but with Go callbacks instead of C callbacks.
var inputST struct {
	createCB  Callback // User-provided callback for creating work
	processCB Callback // User-provided callback for processing a queue
}

// CallbackCreate specifies a user-provided callback that will enqueue
// work when asked.
func CallbackCreate(cb Callback) {
	inputST.createCB = cb
	C.invoke_cb_create()
}

// CallbackProcess specifies a user-provided callback that will
// dequeue and perform work when asked.  Note that the callback is
// allowed to call Enqueue to enqueue additional work if desired.
func CallbackProcess(cb Callback) {
	inputST.processCB = cb
	C.invoke_cb_process()
}

// Begin creates and executes work based on the user-provided callback
// functions.
func Begin() {
	C.CIRCLE_begin()
}

// Checkpoint makes each rank dump a checkpoint file of the form
// "circle<rank>.txt".
func Checkpoint() {
	C.CIRCLE_checkpoint()
}

// ReadRestarts initializes the libcircle queues from the restart
// files produced by the Checkpoint function.
func ReadRestarts() {
	C.CIRCLE_read_restarts()
}

// Abort makes each rank dump a checkpoint file (a la the Checkpoint
// function) and exit.
func Abort() {
	C.CIRCLE_abort()
}

// A LogLevel specifies how verbose libcircle should be while it runs.
type LogLevel uint32

// These constants define the various LogLevel values.
const (
	LogFatal = LogLevel(iota) // Output only fatal errors.
	LogErr                    // Output the above plus nonfatal errors.
	LogWarn                   // Output all of the above plus warnings.
	LogInfo                   // Output all of the above plus informational messages.
	LogDbg                    // Output all of the above plus internal debug messages.
)

// EnableLogging sets libcircle's output verbosity.
func EnableLogging(ll LogLevel) {
	C.CIRCLE_enable_logging(uint32(ll))
}

// Wtime returns the time in seconds from an unspecified epoch.  It
// can be used for benchmarking purposes (although it's a bit
// redundant with Go's time package).
func Wtime() float64 {
	return float64(C.CIRCLE_wtime())
}

// Finalize shuts down libcircle.
func Finalize() {
	C.CIRCLE_finalize()
}
