/*
The following functions provide Go language callbacks that libcircle
can invoke.  They wrap libcircle datatypes with Go interfaces for a
more Go-like feel.  Because the functions are declared with "//export"
they are incompatible with cgo files that declare C code, such as
circle.go.  Therefore, callUserCBCreate and callUserCBProcess must
appear in a separate file.
*/

package circle

// #include <libcircle.h>
import "C"

// Convert a libcircle CIRCLE_handle to something that implements the
// Handle interface, and invoke the user-provided "create" callback
// function.
//
//export callUserCBCreate
func callUserCBCreate(hnd *C.CIRCLE_handle) {
	inputST.createCB(*&queue{handle: hnd})
}

// Convert a libcircle CIRCLE_handle to something that implements the
// Handle interface, and invoke the user-provided "process" callback
// function.
//
//export callUserCBProcess
func callUserCBProcess(hnd *C.CIRCLE_handle) {
	inputST.processCB(*&queue{handle: hnd})
}
