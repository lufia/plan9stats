// Package stats provides statistic utilities for Plan 9.
package stats

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Host represents host status.
type Host struct {
	CPU      *CPU
	Storages []*Storage
}

type CPU struct {
	Name string
	MHz  int // clock rate in MHz
}

type Storage struct {
	Name     string
	Model    string
	Capacity int64
}

// ReadHost returns host status.
func ReadHost(rootdir string) (*Host, error) {
	cpu, err := readCPUType(rootdir)
	if err != nil {
		return nil, err
	}
	var h Host
	h.CPU = cpu

	a, err := readStorages(rootdir)
	if err != nil {
		return nil, err
	}
	h.Storages = a
	return &h, nil
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
		Name: string(b[:i]),
		MHz:  clock,
	}, nil
}

var (
	delim = []byte{' '}
)

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
			sec, err := strconv.ParseInt(string(fields[1]), 10, 64)
			if err != nil {
				return nil, err
			}
			size, err := strconv.ParseInt(string(fields[2]), 10, 64)
			if err != nil {
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
