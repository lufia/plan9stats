package stats

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestReadSysstat(t *testing.T) {
	stat, err := ReadSysstat("testdata")
	if err != nil {
		t.Fatal(err)
	}
	want := []*Sysstat{
		&Sysstat{
			ID:           0,
			NumCtxSwitch: 59251106,
			NumInterrupt: 37524162,
			NumSyscall:   1208203,
			NumFault:     65907,
			NumTLBFault:  0,
			NumTLBPurge:  0,
			LoadAvg:      7,
			Idle:         100,
			Interrupt:    0,
		},
		&Sysstat{
			ID:           1,
			NumCtxSwitch: 219155408,
			NumInterrupt: 28582838,
			NumSyscall:   5017097,
			NumFault:     1002072,
			NumTLBFault:  0,
			NumTLBPurge:  0,
			LoadAvg:      0,
			Idle:         98,
			Interrupt:    1,
		},
	}
	if !cmp.Equal(stat, want) {
		t.Errorf("ReadSysstat: %v", cmp.Diff(stat, want))
	}
}

func TestReadIfaceStats(t *testing.T) {
	stats, err := ReadIfaceStats("testdata/net/ether0")
	if err != nil {
		t.Fatal(err)
	}
	want := &IfaceStats{
		In:               11645518,
		Link:             0,
		Out:              269378,
		NumCRCErr:        0,
		NumOverflows:     0,
		NumSoftOverflows: 0,
		NumFramingErr:    0,
		NumBufferingErr:  0,
		NumOutputErr:     0,
		Promiscuous:      0,
		Mbps:             1000,
		Addr:             "525409008379",
	}
	if !cmp.Equal(stats, want) {
		t.Errorf("ReadIfaceStats: %v", cmp.Diff(stats, want))
	}
}
