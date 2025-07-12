package monitor

import (
	"io/ioutil"
	"log/slog"
	"runtime"
	"strconv"
	"strings"
	"sync/atomic"
	"syscall"
)

var reqCount int64

// CheckLimitsEveryN increments the request counter and, every N requests, logs OS resource limits
func CheckLimitsEveryN(n int64) {
	c := atomic.AddInt64(&reqCount, 1)
	if c%n == 0 {
		logLimits()
	}
}

func logLimits() {
	// RLIMIT_NOFILE (max open files)
	var rLimit syscall.Rlimit
	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit); err == nil {
		slog.Info("[MONITOR] RLIMIT_NOFILE", "cur", rLimit.Cur, "max", rLimit.Max)
	} else {
		slog.Warn("[MONITOR] Could not get RLIMIT_NOFILE", "error", err)
	}

	// RLIMIT_NPROC (max processes/threads)
	logNprocLimit()

	// Go runtime stats
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	slog.Info("[MONITOR] Go MemStats",
		"NumGoroutine", runtime.NumGoroutine(),
		"AllocMB", m.Alloc/1024/1024,
		"SysMB", m.Sys/1024/1024,
		"HeapAllocMB", m.HeapAlloc/1024/1024,
		"HeapSysMB", m.HeapSys/1024/1024,
		"HeapObjects", m.HeapObjects,
		"StackInuseMB", m.StackInuse/1024/1024,
		"StackSysMB", m.StackSys/1024/1024,
		"NumGC", m.NumGC,
		"PauseTotalNs", m.PauseTotalNs,
		"LastGC", m.LastGC,
	)

	// CPUs
	slog.Info("[MONITOR] CPUs", "NumCPU", runtime.NumCPU())

	// Uptime
	uptime := getUptimeSeconds()
	slog.Info("[MONITOR] Uptime", "seconds", uptime)

	// Open file descriptors (Linux/macOS only)
	if numFds, err := countOpenFDs(); err == nil {
		slog.Info("[MONITOR] OpenFDs", "count", numFds)
	}

	// TCP sockets (Linux/macOS only)
	if tcpCount, err := countTCPSockets(); err == nil {
		slog.Info("[MONITOR] TCP Sockets", "count", tcpCount)
	}

	// Page faults (Linux/macOS only)
	if pfaults, err := getPageFaults(); err == nil {
		slog.Info("[MONITOR] PageFaults", "minor", pfaults.Minor, "major", pfaults.Major)
	}
}

type pageFaults struct{ Minor, Major uint64 }

func getUptimeSeconds() int64 {
	if data, err := ioutil.ReadFile("/proc/uptime"); err == nil {
		parts := strings.Fields(string(data))
		if len(parts) > 0 {
			if secs, err := strconv.ParseFloat(parts[0], 64); err == nil {
				return int64(secs)
			}
		}
	}
	return -1
}

func countOpenFDs() (int, error) {
	fdDir := "/proc/self/fd"
	files, err := ioutil.ReadDir(fdDir)
	if err != nil {
		return -1, err
	}
	return len(files), nil
}

func countTCPSockets() (int, error) {
	// Conta linhas em /proc/net/tcp
	if data, err := ioutil.ReadFile("/proc/net/tcp"); err == nil {
		return len(strings.Split(strings.TrimSpace(string(data)), "\n")) - 1, nil
	}
	return -1, nil
}

func getPageFaults() (pageFaults, error) {
	pf := pageFaults{}
	if data, err := ioutil.ReadFile("/proc/self/stat"); err == nil {
		fields := strings.Fields(string(data))
		if len(fields) > 12 {
			// minflt: 10, majflt: 12
			pf.Minor, _ = strconv.ParseUint(fields[9], 10, 64)
			pf.Major, _ = strconv.ParseUint(fields[11], 10, 64)
		}
		return pf, nil
	}
	return pf, nil
}
