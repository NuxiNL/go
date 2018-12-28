// +build cloudabi

package syscall

type Errno uint16

func (e Errno) Error() string {
	return "errno " + itoa(int(e))
}

type SysProcAttr struct {}
