// +build cloudabi

package runtime

const (
	active_spin     = 4
	active_spin_cnt = 30
)

func lock(l *mutex) {
}

func unlock(l *mutex) {
}

func noteclear(n *note) {}

func notewakeup(n *note) {}

func notetsleep(n *note, ns int64) bool {
	return false
}

func notetsleepg(n *note, ns int64) bool {
	return false
}

func notesleep(n *note) {}

func checkTimeouts() {
}

func beforeIdle() bool {
	return false
}
