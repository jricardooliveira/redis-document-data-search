//go:build darwin

package monitor

func logNprocLimit() {
	// RLIMIT_NPROC not available on macOS
}
