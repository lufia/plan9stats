package stats

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestReadSysStats(t *testing.T) {
	ctx := context.Background()
	stat, err := ReadSysStats(ctx, WithRootDir("testdata"))
	if err != nil {
		t.Fatal(err)
	}
	want := []*SysStats{
		&SysStats{
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
		&SysStats{
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
	if !cmp.Equal(want, stat) {
		t.Errorf("ReadSysStats: %v", cmp.Diff(want, stat))
	}
}

func TestReadInterfaceStats(t *testing.T) {
	ctx := context.Background()
	stats, err := ReadInterfaceStats(ctx, WithRootDir("testdata/net/ether0"))
	if err != nil {
		t.Fatal(err)
	}
	want := &InterfaceStats{
		PacketsReceived:  11645518,
		Link:             0,
		PacketsSent:      269378,
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
	if !cmp.Equal(want, stats) {
		t.Errorf("ReadInterfaceStats: %v", cmp.Diff(want, stats))
	}
}
