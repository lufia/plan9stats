package stats

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestReadHost(t *testing.T) {
	h, err := ReadHost("testdata")
	if err != nil {
		t.Fatal(err)
	}
	want := &Host{
		CPU: &CPU{
			Name:  "Core i7/Xeon",
			Clock: 2403,
		},
		Memory: &Memory{
			Total:        1071185920,
			PageSize:     4096,
			KernelPages:  61372,
			UserPages:    Gauge{Used: 2792, Avail: 200148},
			SwapPages:    Gauge{Used: 0, Avail: 160000},
			KernelMalloc: Gauge{Used: 9046176, Avail: 219352384},
			KernelDraw:   Gauge{Used: 0, Avail: 16777216},
		},
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
	}
	if !cmp.Equal(h, want) {
		t.Errorf("ReadHost: %v", cmp.Diff(h, want))
	}
}
