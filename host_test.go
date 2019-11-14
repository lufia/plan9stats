package stats

import (
	"testing"
)

func TestReadHost(t *testing.T) {
	h, err := ReadHost("testdata")
	if err != nil {
		t.Fatal(err)
	}
	if want := "Core i7/Xeon"; h.CPU.Name != want {
		t.Errorf("CPU.Name = %v; want %v", h.CPU.Name, want)
	}
	if want := 2403; h.CPU.MHz != want {
		t.Errorf("CPU.MHz = %v; want %v", h.CPU.MHz, want)
	}

	if want := 2; len(h.Storages) != want {
		t.Errorf("len(Storages) = %d; want %d", len(h.Storages), want)
	}
}
