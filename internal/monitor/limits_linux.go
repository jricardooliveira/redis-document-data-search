//go:build linux

package monitor

import (
	"log/slog"
	"syscall"
)

func logNprocLimit() {
	var pLimit syscall.Rlimit
	if err := syscall.Getrlimit(syscall.RLIMIT_NPROC, &pLimit); err == nil {
		slog.Info("[MONITOR] RLIMIT_NPROC", "cur", pLimit.Cur, "max", pLimit.Max)
	}
}
