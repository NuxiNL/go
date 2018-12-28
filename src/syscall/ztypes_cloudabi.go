// +build cloudabi

package syscall

type Timespec struct {
	Sec  int64
	Nsec int32
}

type Timeval struct {
	Sec  int64
	Usec int32
}

type Stat_t struct {
	Dev      uint64
	Ino      uint64
}
