Circle
======

[![Go Reference](https://pkg.go.dev/badge/github.com/lanl/circle.svg)](https://pkg.go.dev/github.com/lanl/circle) [![Go Report Card](https://goreportcard.com/badge/github.com/lanl/circle)](https://goreportcard.com/report/github.com/lanl/circle)

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

You'll need to download and install Libcircle, which is available from <https://github.com/hpc/libcircle>.  After that,
```bash
go mod tidy
```
ought to work for any program that includes an
```Go
import "github.com/lanl/circle"
```

Documentation
-------------

Pre-built documentation for the Circle API is available online at <https://pkg.go.dev/github.com/lanl/circle>.

Legal statement
---------------

Copyright © 2011, Triad National Security, LLC
All rights reserved.

This software was produced under U.S. Government contract 89233218CNA000001 for Los Alamos National Laboratory (LANL), which is operated by Triad National Security, LLC for the U.S. Department of Energy/National Nuclear Security Administration.  All rights in the program are reserved by Triad National Security, LLC, and the U.S. Department of Energy/National Nuclear Security Administration. The Government is granted for itself and others acting on its behalf a nonexclusive, paid-up, irrevocable worldwide license in this material to reproduce, prepare derivative works, distribute copies to the public, perform publicly and display publicly, and to permit others to do so.  NEITHER THE GOVERNMENT NOR TRIAD NATIONAL SECURITY, LLC MAKES ANY WARRANTY, EXPRESS OR IMPLIED, OR ASSUMES ANY LIABILITY FOR THE USE OF THIS SOFTWARE.  If software is modified to produce derivative works, such modified software should be clearly marked, so as not to confuse it with the version available from LANL.

Circle is provided under a BSD-ish license with a "modifications must be indicated" clause.  See [the LICENSE file](http://github.com/lanl/circle/blob/master/LICENSE.md) for the full text.

Circle is part of the LANL Go Suite, identified internally by LANL as LA-CC-11-056.

Author
------

Scott Pakin, <pakin@lanl.gov>
