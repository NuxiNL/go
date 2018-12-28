// +build cloudabi

package runtime

func initsig(preinit bool) {}

func msigsave(mp *m) {}

func sigblock() {}

func signame(sig uint32) string {
	return ""
}

func crash() {
}

func msigrestore(sigmask sigset) {}
