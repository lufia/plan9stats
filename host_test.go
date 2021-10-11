package stats

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestReadHost(t *testing.T) {
	ctx := context.Background()
	h, err := ReadHost(ctx, WithRootDir("testdata"))
	if err != nil {
		t.Fatal(err)
	}
	want := &Host{
		Sysname: "gnot",
		Storages: []*Storage{
			&Storage{
				Name:     "sdC0",
				Model:    "QEMU HARDDISK",
				Capacity: 209715200 * 512, // 100GB
			},
			&Storage{
				Name:  "sdD0",
				Model: "QEMU    QEMU DVD-ROM    0.12",
			},
		},
		Interfaces: []*Interface{
			&Interface{
				Name: "ether0",
				Addr: "525409008379",
			},
		},
	}
	if !cmp.Equal(want, h) {
		t.Errorf("ReadHost: %v", cmp.Diff(want, h))
	}
}

func TestReadMemStats(t *testing.T) {
	ctx := context.Background()
	h, err := ReadMemStats(ctx, WithRootDir("testdata"))
	if err != nil {
		t.Fatal(err)
	}
	want := &MemStats{
		Total:       1071185920,
		PageSize:    4096,
		KernelPages: 61372,
		UserPages:   Gauge{Used: 2792, Avail: 200148},
		SwapPages:   Gauge{Used: 0, Avail: 160000},
		Malloced:    Gauge{Used: 9046176, Avail: 219352384},
		Graphics:    Gauge{Used: 0, Avail: 16777216},
	}
	if !cmp.Equal(want, h) {
		t.Errorf("ReadMemStats: %v", cmp.Diff(want, h))
	}
}
