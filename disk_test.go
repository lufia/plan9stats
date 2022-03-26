package stats

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestReadStorages(t *testing.T) {
	ctx := context.Background()
	disks, err := ReadStorages(ctx, WithRootDir("testdata"))
	if err != nil {
		t.Fatal(err)
	}
	want := []*Storage{
		&Storage{
			Name:     "sdC0",
			Model:    "QEMU HARDDISK",
			Capacity: 209715200 * 512, // 100GB
			Partitions: []*Partition{
				{"data", 0, 209715200},
				{"plan9", 63, 209712510},
				{"9fat", 63, 204863},
				{"nvram", 204863, 204864},
				{"fossil", 204864, 208663934},
				{"swap", 208663934, 209712510},
			},
		},
		&Storage{
			Name:  "sdD0",
			Model: "QEMU    QEMU DVD-ROM    0.12",
		},
	}
	if !cmp.Equal(want, disks) {
		t.Errorf("ReadHost: %v", cmp.Diff(want, disks))
	}
}
