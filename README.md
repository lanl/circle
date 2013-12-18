Circle
======

Description
-----------

The Circle package provides a Go interface to the [Libcircle](http://hpc.github.io/libcircle/) distributed-queue API.  Despite the name, Circle has nothing to do with graphics.  Instead, Circle provides a mechanism for enqueueing "work" (currently, text strings) on a distributed queue then letting numerous processes distributed across a local-area network dequeue and process that work.

Use Circle when you have a huge number of independent tasks to perform and want an easy way to distribute these across a large cluster or supercomputer.

Features
--------

The underlying Libcircle library has the following features:

* proximity-aware, work-stealing scheduler

* used daily on production supercomputers at [Los Alamos National Laboratory](http://www.lanl.gov/) to perform various maintenance activities across a multi-petabyte parallel filesystem

* fast&mdash;communication is implemented with user-level messaging (specifically, [MPI](http://www.mpi-forum.org/)), not kernel-level sockets.

Circle provides a Go interface to Libcircle:

* a low-level API that maps directly to the Libcircle API but supports all of the Go niceties such as using Go strings for work items and Go functions for Libcircle callbacks

* a higher-level API that forgoes Libcircle's callback mechanism in favor Go channels: one for enqueueing work and one for dequeueing work

Installation
------------

You'll need to download and install Libcircle, which is available from https://github.com/hpc/libcircle.  After that,

    go get github.com/losalamos/circle

ought to work.

License
-------

Circle is provided under a BSD-ish license with a "modifications must be indicated" clause.  See [the LICENSE file](http://github.com/losalamos/circle/blob/master/LICENSE.md) for the full text.

Circle is part of the [LANL Go Suite](http://www.lanl.gov/projects/feynman-center/technologies/software/lanl%20go%20suite.php), LA-CC-11-056.

Author
------

Scott Pakin, <pakin@lanl.gov>
