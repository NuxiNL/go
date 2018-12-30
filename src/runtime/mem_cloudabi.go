// +build cloudabi

package runtime

import (
	"internal/syscall/cloudabi"
	"unsafe"
)

func sysAlloc(n uintptr, sysStat *uint64) unsafe.Pointer {
	v, err := cloudabi_sys_mem_map(
		nil, uint(n), cloudabi.Mprot_Read | cloudabi.Mprot_Write,
		cloudabi.Mflags_Anon | cloudabi.Mflags_Private,
		cloudabi.Fd_MapAnonFd, 0)
	if err != cloudabi.Errno_Success {
		return nil
	}
	mSysStatInc(sysStat, n)
	return v
}

func sysUnused(v unsafe.Pointer, n uintptr) {
	cloudabi_sys_mem_advise(v, uint(n), cloudabi.Advice_Dontneed)
}

func sysUsed(v unsafe.Pointer, n uintptr) {
}

// Don't split the stack as this function may be invoked without a valid G,
// which prevents us from allocating more stack.
//go:nosplit
func sysFree(v unsafe.Pointer, n uintptr, sysStat *uint64) {
	mSysStatDec(sysStat, n)
	cloudabi_sys_mem_unmap(v, uint(n))
}

func sysFault(v unsafe.Pointer, n uintptr) {
	cloudabi_sys_mem_map(
		v, uint(n), 0, cloudabi.Mflags_Anon | cloudabi.Mflags_Fixed |
		cloudabi.Mflags_Private, cloudabi.Fd_MapAnonFd, 0)
}

func sysReserve(v unsafe.Pointer, n uintptr) unsafe.Pointer {
	p, err := cloudabi_sys_mem_map(
		v, uint(n), 0, cloudabi.Mflags_Anon | cloudabi.Mflags_Private,
		cloudabi.Fd_MapAnonFd, 0)
	if err != cloudabi.Errno_Success {
		return nil
	}
	return p
}

func sysMap(v unsafe.Pointer, n uintptr, sysStat *uint64) {
	mSysStatInc(sysStat, n)

	p, err := cloudabi_sys_mem_map(
		v, uint(n), cloudabi.Mprot_Read | cloudabi.Mprot_Write,
		cloudabi.Mflags_Anon | cloudabi.Mflags_Fixed |
		cloudabi.Mflags_Private, cloudabi.Fd_MapAnonFd, 0)
	if err == cloudabi.Errno_Nomem {
		throw("runtime: out of memory")
	}
	if p != v || err != cloudabi.Errno_Success {
		throw("runtime: cannot map pages in arena address space")
	}
}
