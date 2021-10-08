// Package stats provides statistic utilities for Plan 9.
package stats

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var (
	delim = []byte{' '}
)

// Host represents host status.
type Host struct {
	Sysname    string
	CPU        *CPU
	Storages   []*Storage
	Interfaces []*Interface
}

// CPU
type CPU struct {
	Name  string
	Clock int // clock rate in MHz
}

// MemStats represents the memory statistics.
type MemStats struct {
	Total       int64 // total memory in byte
	PageSize    int64 // a page size in byte
	KernelPages int64
	UserPages   Gauge
	SwapPages   Gauge

	Malloced Gauge // kernel malloced data in byte
	Graphics Gauge // kernel graphics data in byte
}

// Gauge is used/available gauge.
type Gauge struct {
	Used  int64
	Avail int64
}

func (g Gauge) Free() int64 {
	return g.Avail - g.Used
}

// ReadMemStats reads memory statistics from /dev/swap.
func ReadMemStats(ctx context.Context, opts ...Option) (*MemStats, error) {
	cfg := newConfig(opts...)
	swap := filepath.Join(cfg.rootdir, "/dev/swap")
	f, err := os.Open(swap)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var stat MemStats
	m := map[string]interface{}{
		"memory":        &stat.Total,
		"pagesize":      &stat.PageSize,
		"kernel":        &stat.KernelPages,
		"user":          &stat.UserPages,
		"swap":          &stat.SwapPages,
		"kernel malloc": &stat.Malloced,
		"kernel draw":   &stat.Graphics,
	}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		fields := bytes.SplitN(scanner.Bytes(), delim, 2)
		if len(fields) < 2 {
			continue
		}
		switch key := string(fields[1]); key {
		case "memory", "pagesize", "kernel":
			v := m[key].(*int64)
			n, err := strconv.ParseInt(string(fields[0]), 10, 64)
			if err != nil {
				return nil, err
			}
			*v = n
		case "user", "swap", "kernel malloc", "kernel draw":
			v := m[key].(*Gauge)
			if err := parseGauge(string(fields[0]), v); err != nil {
				return nil, err
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return &stat, nil
}

func parseGauge(s string, r *Gauge) error {
	a := strings.SplitN(s, "/", 2)
	if len(a) != 2 {
		return fmt.Errorf("can't parse ratio: %s", s)
	}
	var p intParser
	u := p.ParseInt64(a[0], 10)
	n := p.ParseInt64(a[1], 10)
	if err := p.Err(); err != nil {
		return err
	}
	r.Used = u
	r.Avail = n
	return nil
}

type Storage struct {
	Name     string
	Model    string
	Capacity int64
}

type Interface struct {
	Name string
	Addr string
}

const (
	numEther = 8  // see ether(3)
	numIpifc = 16 // see ip(3)
)

// ReadInterfaces reads network interfaces from etherN.
func ReadInterfaces(ctx context.Context, opts ...Option) ([]*Interface, error) {
	cfg := newConfig(opts...)
	var a []*Interface
	for i := 0; i < numEther; i++ {
		p, err := readInterface(cfg.rootdir, i)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return nil, err
		}
		a = append(a, p)
	}
	return a, nil
}

func readInterface(netroot string, i int) (*Interface, error) {
	ether := fmt.Sprintf("ether%d", i)
	dir := filepath.Join(netroot, ether)
	info, err := os.Stat(dir)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("%s: is not directory", dir)
	}

	addr, err := ioutil.ReadFile(filepath.Join(dir, "addr"))
	if err != nil {
		return nil, err
	}
	return &Interface{
		Name: ether,
		Addr: string(addr),
	}, nil
}

var (
	netdirs = []string{"/net", "/net.alt"}
)

// ReadHost reads host status.
func ReadHost(ctx context.Context, opts ...Option) (*Host, error) {
	cfg := newConfig(opts...)
	var h Host
	name, err := readSysname(cfg.rootdir)
	if err != nil {
		return nil, err
	}
	h.Sysname = name

	cpu, err := readCPUType(cfg.rootdir)
	if err != nil {
		return nil, err
	}
	h.CPU = cpu

	a, err := readStorages(cfg.rootdir)
	if err != nil {
		return nil, err
	}
	h.Storages = a

	for _, s := range netdirs {
		netroot := filepath.Join(cfg.rootdir, s)
		ifaces, err := ReadInterfaces(ctx, WithRootDir(netroot))
		if err != nil {
			return nil, err
		}
		h.Interfaces = append(h.Interfaces, ifaces...)
	}
	return &h, nil
}

func readSysname(rootdir string) (string, error) {
	file := filepath.Join(rootdir, "/dev/sysname")
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}
	return string(bytes.TrimSpace(b)), nil
}

func readCPUType(rootdir string) (*CPU, error) {
	file := filepath.Join(rootdir, "/dev/cputype")
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	b = bytes.TrimSpace(b)
	i := bytes.LastIndexByte(b, ' ')
	if i < 0 {
		return nil, fmt.Errorf("%s: invalid format", file)
	}
	clock, err := strconv.Atoi(string(b[i+1:]))
	if err != nil {
		return nil, err
	}
	return &CPU{
		Name:  string(b[:i]),
		Clock: clock,
	}, nil
}

func readStorages(rootdir string) ([]*Storage, error) {
	sdctl := filepath.Join(rootdir, "/dev/sdctl")
	f, err := os.Open(sdctl)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var a []*Storage
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		fields := bytes.Split(scanner.Bytes(), delim)
		if len(fields) == 0 {
			continue
		}
		exp := string(fields[0]) + "*"
		if !strings.HasPrefix(exp, "sd") {
			continue
		}
		dir := filepath.Join(rootdir, "/dev", exp)
		m, err := filepath.Glob(dir)
		if err != nil {
			return nil, err
		}
		for _, dir := range m {
			s, err := readStorage(dir)
			if err != nil {
				return nil, err
			}
			a = append(a, s)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return a, nil
}

func readStorage(dir string) (*Storage, error) {
	ctl := filepath.Join(dir, "ctl")
	f, err := os.Open(ctl)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var s Storage
	s.Name = filepath.Base(dir)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Bytes()
		switch {
		case bytes.HasPrefix(line, []byte("inquiry")):
			s.Model = string(bytes.TrimSpace(line[7:]))
		case bytes.HasPrefix(line, []byte("geometry")):
			fields := bytes.Split(line, delim)
			if len(fields) < 3 {
				continue
			}
			var p intParser
			sec := p.ParseInt64(string(fields[1]), 10)
			size := p.ParseInt64(string(fields[2]), 10)
			if err := p.Err(); err != nil {
				return nil, err
			}
			s.Capacity = sec * size
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return &s, nil
}

type IPStats struct {
	ID      int    // number of interface in ipifc dir
	Device  string // associated physical device
	MTU     int    // max transfer unit
	Sendra6 uint8  // on == send router adv
	Recvra6 uint8  // on == recv router adv

	Pktin  int64 // packets read
	Pktout int64 // packets written
	Errin  int64 // read errors
	Errout int64 // write errors
}

type Iplifc struct {
	IP            net.IP
	Mask          net.IPMask
	Net           net.IP // ip & mask
	PerfLifetime  int64  // preferred lifetime
	ValidLifetime int64  // valid lifetime
}

type Ipv6rp struct {
	// TODO(lufia): see ip(2)
}
