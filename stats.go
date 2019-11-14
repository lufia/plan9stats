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
		fields := bytes.Split(scanner.Bytes(), []byte{' '})
		_ = fields[0] // TODO(lufia): read subdirs
		a = append(a, &Storage{})
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return a, nil
}
