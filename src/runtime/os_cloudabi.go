// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build cloudabi

package runtime

import (
	"internal/syscall/cloudabi"
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

func sysargs(argc int32, argv **byte) {
	sysauxv((*[1 << 20]cloudabi.Auxv)(unsafe.Pointer(argv))[:])
}

var (
	_cloudabi_sys_mem_map unsafe.Pointer
)

func linkVdso(ehdr *elfEhdr) {
	// Extract the Dynamic Section of the vDSO.
	base := unsafe.Pointer(ehdr)
	phdrList := (*[1 << 20]elfPhdr)(add(base, uintptr(ehdr.e_phoff)))[:ehdr.e_phnum]
	var dynList *[1 << 20]elfDyn
	for _, phdr := range phdrList {
		switch phdr.p_type {
		case _PT_DYNAMIC:
			dynList = (*[1 << 20]elfDyn)(add(base, uintptr(phdr.p_offset)))
		}
	}

	// Extract the symbol and string tables.
	var symsz uint32
	var str unsafe.Pointer
	var symBase *[1 << 20]elfSym
	for _, dyn := range dynList {
		if dyn.d_tag == _DT_NULL {
			break
		}
		switch dyn.d_tag {
		case _DT_HASH:
			// Number of symbols in the symbol table can
			// only be extracted by fetching the number of
			// chains in the symbol hash table.
			symsz = (*[2]uint32)(add(base, uintptr(dyn.d_val)))[1]
		case _DT_STRTAB:
			str = add(base, uintptr(dyn.d_val))
		case _DT_SYMTAB:
			symBase = (*[1 << 20]elfSym)(add(base, uintptr(dyn.d_val)))
		}
	}

	// Scan through all of the symbols and find the implementations
	// of the system calls.
	symList := symBase[:symsz]
	for _, sym := range symList {
		switch gostringnocopy((*byte)(add(str, uintptr(sym.st_name)))) {
		case "cloudabi_sys_mem_map":
			_cloudabi_sys_mem_map = add(base, uintptr(sym.st_value))
		}
	}
}

func sysauxv(auxv []cloudabi.Auxv) {
	for i := 0; auxv[i].AType != cloudabi.Auxtype_Null; i++ {
		switch auxv[i].AType {
		case cloudabi.Auxtype_Pagesz:
			physPageSize = *(*uintptr)(unsafe.Pointer(&auxv[i].AuxvUnion))
		case cloudabi.Auxtype_SysinfoEhdr:
			linkVdso(*(**elfEhdr)(unsafe.Pointer(&auxv[i].AuxvUnion)))
		}
	}
}
