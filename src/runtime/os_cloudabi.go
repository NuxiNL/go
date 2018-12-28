// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build cloudabi

package runtime

import (
	"unsafe"
)

const _NSIG = 27

func getRandomData(r []byte) {}

type sigset struct{}

type mOS struct{}

func osyield() {}

type gsignalStack struct{}

func mpreinit(mp *m) {
	mp.gsignal = malg(32 * 1024)
	mp.gsignal.m = mp
}

func minit() {
}

func unminit() {
}

func newosproc(mp *m) {
	panic("newosproc: not implemented")
}

func goenvs() {}

func setProcessCPUProfiler(hz int32) {}
func setThreadCPUProfiler(hz int32)  {}
func sigdisable(uint32)              {}
func sigenable(uint32)               {}
func sigignore(uint32)               {}
func clearSignalHandlers()           {}

func nanotime() int64 { return 123 }

func write(fd uintptr, p unsafe.Pointer, n int32) int32 {
	return n
}

func exit(code int32) {}

func usleep(usec uint32) {}

func exitThread(wait *uint32) {}

func walltime() (sec int64, nsec int32) {
	return 0, 0
}

func settls() {}
func osinit() {}
