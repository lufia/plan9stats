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
			Name: "Core i7/Xeon",
			MHz:  2403,
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
