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
		Sysname: "gnot",
		CPU: &CPU{
			Name:  "Core i7/Xeon",
			Clock: 2403,
		},
		MemStats: &MemStats{
			Total:       1071185920,
			PageSize:    4096,
			KernelPages: 61372,
			UserPages:   Ratio{Used: 2792, Avail: 200148},
			SwapPages:   Ratio{Used: 0, Avail: 160000},
			Malloced:    Ratio{Used: 9046176, Avail: 219352384},
			Graphics:    Ratio{Used: 0, Avail: 16777216},
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
		Interfaces: []*Interface{
			&Interface{
				Name: "ether0",
				Addr: "525409008379",
			},
		},
	}
	if !cmp.Equal(h, want) {
		t.Errorf("ReadHost: %v", cmp.Diff(h, want))
	}
}
