package stats

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestReadCPUStats(t *testing.T) {
	ctx := context.Background()
	stat, err := ReadCPUStats(ctx, WithRootDir("testdata"))
	if err != nil {
		t.Fatal(err)
	}
	const (
		userTime = (10 + 2380 + 0 + 390) * time.Millisecond
		sysTime  = (20 + 29690 + 0 + 310) * time.Millisecond
	)
	want := &CPUStats{
		User: userTime,
		Sys:  sysTime,
		Idle: (1412961713341830*2 - userTime - sysTime),
	}
	if !cmp.Equal(want, stat) {
		t.Errorf("ReadCPUTime: %v", cmp.Diff(want, stat))
	}
}

func TestReadTime(t *testing.T) {
	ctx := context.Background()
	stat, err := ReadTime(ctx, WithRootDir("testdata"))
	if err != nil {
		t.Fatal(err)
	}
	want := &Time{
		Unix:     1633882064 * time.Second,
		UnixNano: 1633882064926300833 * time.Nanosecond,
		Ticks:    2825920097745864,
		Freq:     1999997644,
	}
	if !cmp.Equal(want, stat) {
		t.Errorf("ReadCPUTime: %v", cmp.Diff(want, stat))
	}
}
