// +build cloudabi

package runtime

import (
	"internal/syscall/cloudabi"
	"unsafe"
)

func sysAlloc(n uintptr, sysStat *uint64) unsafe.Pointer {
	// TODO: This is incorrect.
	var v unsafe.Pointer
	args := [...]unsafe.Pointer{
		nil,
		unsafe.Pointer(uintptr(n)),
		unsafe.Pointer(uintptr(cloudabi.Mprot_Read | cloudabi.Mprot_Write)),
		unsafe.Pointer(uintptr(cloudabi.Mflags_Anon | cloudabi.Mflags_Private)),
		unsafe.Pointer(uintptr(cloudabi.Fd_MapAnonFd)),
		nil,
		unsafe.Pointer(&v),
	}
	err := asmcgocall(_cloudabi_sys_mem_map, unsafe.Pointer(&args))
	if err != 0 {
		return nil
	}
	mSysStatInc(sysStat, n)
	return v
}

func sysMap(v unsafe.Pointer, n uintptr, sysStat *uint64) {
}

func sysUnused(v unsafe.Pointer, n uintptr) {
}

func sysUsed(v unsafe.Pointer, n uintptr) {
}

func sysFree(v unsafe.Pointer, n uintptr, sysStat *uint64) {
}

func sysFault(v unsafe.Pointer, n uintptr) {
}

func sysReserve(v unsafe.Pointer, n uintptr) unsafe.Pointer {
	return nil
}
