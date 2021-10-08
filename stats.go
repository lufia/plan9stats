package stats

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

type SysStats struct {
	ID           int
	NumCtxSwitch int64
	NumInterrupt int64
	NumSyscall   int64
	NumFault     int64
	NumTLBFault  int64
	NumTLBPurge  int64
	LoadAvg      int64 // in units of milli-CPUs and is decayed over time
	Idle         int   // percentage
	Interrupt    int   // percentage
}

// ReadSysStats reads system statistics from /dev/sysstat.
func ReadSysStats(opts ...Option) ([]*SysStats, error) {
	cfg := newConfig(opts...)
	file := filepath.Join(cfg.rootdir, "/dev/sysstat")
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var stats []*SysStats
	for scanner.Scan() {
		a := strings.Fields(scanner.Text())
		if len(a) != 10 {
			continue
		}
		var (
			p    intParser
			stat SysStats
		)
		stat.ID = p.ParseInt(a[0], 10)
		stat.NumCtxSwitch = p.ParseInt64(a[1], 10)
		stat.NumInterrupt = p.ParseInt64(a[2], 10)
		stat.NumSyscall = p.ParseInt64(a[3], 10)
		stat.NumFault = p.ParseInt64(a[4], 10)
		stat.NumTLBFault = p.ParseInt64(a[5], 10)
		stat.NumTLBPurge = p.ParseInt64(a[6], 10)
		stat.LoadAvg = p.ParseInt64(a[7], 10)
		stat.Idle = p.ParseInt(a[8], 10)
		stat.Interrupt = p.ParseInt(a[9], 10)
		if err := p.Err(); err != nil {
			return nil, err
		}
		stats = append(stats, &stat)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return stats, nil
}

type InterfaceStats struct {
	PacketsReceived  int64 // in packets
	Link             int   // link status
	PacketsSent      int64 // out packets
	NumCRCErr        int   // input CRC errors
	NumOverflows     int   // packet overflows
	NumSoftOverflows int   // software overflow
	NumFramingErr    int   // framing errors
	NumBufferingErr  int   // buffering errors
	NumOutputErr     int   // output errors
	Promiscuous      int   // number of promiscuous opens
	Mbps             int   // megabits per sec
	Addr             string
}

func ReadInterfaceStats(opts ...Option) (*InterfaceStats, error) {
	cfg := newConfig(opts...)
	file := filepath.Join(cfg.rootdir, "stats")
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var stats InterfaceStats
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		s := strings.TrimSpace(scanner.Text())
		a := strings.SplitN(s, ":", 2)
		if len(a) != 2 {
			continue
		}
		var p intParser
		v := strings.TrimSpace(a[1])
		switch a[0] {
		case "in":
			stats.PacketsReceived = p.ParseInt64(v, 10)
		case "link":
			stats.Link = p.ParseInt(v, 10)
		case "out":
			stats.PacketsSent = p.ParseInt64(v, 10)
		case "crc":
			stats.NumCRCErr = p.ParseInt(v, 10)
		case "overflows":
			stats.NumOverflows = p.ParseInt(v, 10)
		case "soft overflows":
			stats.NumSoftOverflows = p.ParseInt(v, 10)
		case "framing errs":
			stats.NumFramingErr = p.ParseInt(v, 10)
		case "buffer errs":
			stats.NumBufferingErr = p.ParseInt(v, 10)
		case "output errs":
			stats.NumOutputErr = p.ParseInt(v, 10)
		case "prom":
			stats.Promiscuous = p.ParseInt(v, 10)
		case "mbps":
			stats.Mbps = p.ParseInt(v, 10)
		case "addr":
			stats.Addr = v
		}
		if err := p.Err(); err != nil {
			return nil, err
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return &stats, nil
}

type TCPStats struct {
	MaxConn            int
	MaxSegment         int
	ActiveOpens        int
	PassiveOpens       int
	EstablishedResets  int
	CurrentEstablished int
}
