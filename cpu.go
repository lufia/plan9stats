package stats

import (
	"context"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Time represents /dev/time.
type Time struct {
	Unix     time.Duration
	UnixNano time.Duration
	Ticks    int64 // clock ticks
	Freq     int64 //cloc frequency
}

// Uptime returns uptime.
func (t *Time) Uptime() time.Duration {
	v := float64(t.Ticks) / float64(t.Freq)
	return time.Duration(v*1000_000_000) * time.Nanosecond
}

func ReadTime(ctx context.Context, opts ...Option) (*Time, error) {
	cfg := newConfig(opts...)
	file := filepath.Join(cfg.rootdir, "/dev/time")
	var t Time
	if err := readTime(file, &t); err != nil {
		return nil, err
	}
	return &t, nil
}

// ProcStatus represents a /proc/n/status.
type ProcStatus struct {
	Name         string
	User         string
	State        string
	Times        CPUTime
	MemUsed      int64  // in units of 1024 bytes
	BasePriority uint32 // 0(low) to 19(high)
	Priority     uint32 // 0(low) to 19(high)
}

// CPUTime represents /dev/cputime or a part of /proc/n/status.
type CPUTime struct {
	User      time.Duration // the time in user mode (millisecconds)
	Sys       time.Duration
	Real      time.Duration
	ChildUser time.Duration // exited children and descendants time in user mode
	ChildSys  time.Duration
	ChildReal time.Duration
}

// CPUStats emulates Linux's /proc/stat.
type CPUStats struct {
	User time.Duration
	Sys  time.Duration
	Idle time.Duration
}

func ReadCPUStats(ctx context.Context, opts ...Option) (*CPUStats, error) {
	cfg := newConfig(opts...)
	dir := filepath.Join(cfg.rootdir, "/proc")
	d, err := os.Open(dir)
	if err != nil {
		return nil, err
	}
	defer d.Close()

	names, err := d.Readdirnames(0)
	if err != nil {
		return nil, err
	}
	var up uint32parser
	pids := make([]uint32, len(names))
	for i, s := range names {
		pids[i] = up.Parse(s)
	}
	if up.err != nil {
		return nil, err
	}
	sort.Slice(pids, func(i, j int) bool {
		return pids[i] < pids[j]
	})

	var stat CPUStats
	for _, pid := range pids {
		s := strconv.FormatUint(uint64(pid), 10)
		file := filepath.Join(dir, s, "status")
		var p ProcStatus
		if err := readProcStatus(file, &p); err != nil {
			return nil, err
		}
		stat.User += p.Times.User
		stat.Sys += p.Times.Sys
	}

	var t Time
	file := filepath.Join(cfg.rootdir, "/dev/time")
	if err := readTime(file, &t); err != nil {
		return nil, err
	}
	stat.Idle = t.Uptime() - stat.User - stat.Sys
	return &stat, nil
}

func readProcStatus(file string, p *ProcStatus) error {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	fields := strings.Fields(string(b))
	if len(fields) != 12 {
		return errors.New("invalid format")
	}
	p.Name = string(fields[0])
	p.User = string(fields[1])
	p.State = string(fields[2])
	var up uint32parser
	p.Times.User = time.Duration(up.Parse(fields[3])) * time.Millisecond
	p.Times.Sys = time.Duration(up.Parse(fields[4])) * time.Millisecond
	p.Times.Real = time.Duration(up.Parse(fields[5])) * time.Millisecond
	p.Times.ChildUser = time.Duration(up.Parse(fields[6])) * time.Millisecond
	p.Times.ChildSys = time.Duration(up.Parse(fields[7])) * time.Millisecond
	p.Times.ChildReal = time.Duration(up.Parse(fields[8])) * time.Millisecond
	p.MemUsed, err = strconv.ParseInt(fields[9], 10, 64)
	if err != nil {
		return err
	}
	p.BasePriority = up.Parse(fields[10])
	p.Priority = up.Parse(fields[11])
	return up.err
}

func readTime(file string, t *Time) error {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	fields := strings.Fields(string(b))
	if len(fields) != 4 {
		return errors.New("invalid format")
	}
	n, err := strconv.ParseInt(fields[0], 10, 32)
	if err != nil {
		return err
	}
	t.Unix = time.Duration(n) * time.Second
	v, err := strconv.ParseInt(fields[1], 10, 64)
	if err != nil {
		return err
	}
	t.UnixNano = time.Duration(v) * time.Nanosecond
	t.Ticks, err = strconv.ParseInt(fields[2], 10, 64)
	if err != nil {
		return err
	}
	t.Freq, err = strconv.ParseInt(fields[3], 10, 64)
	if err != nil {
		return err
	}
	return nil
}

type uint32parser struct {
	err error
}

func (p *uint32parser) Parse(s string) uint32 {
	if p.err != nil {
		return 0
	}
	n, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		p.err = err
		return 0
	}
	return uint32(n)
}
