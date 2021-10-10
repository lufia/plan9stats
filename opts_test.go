package stats

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestOptions(t *testing.T) {
	tests := []struct {
		opts []Option
		cfg  *Config
	}{
		{
			opts: []Option{WithRootDir("/mnt/term")},
			cfg: &Config{
				rootdir: "/mnt/term",
			},
		},
	}
	o := cmp.AllowUnexported(Config{})
	for _, tt := range tests {
		cfg := newConfig(tt.opts...)
		if !cmp.Equal(tt.cfg, cfg, o) {
			t.Errorf("newConfig: %s", cmp.Diff(tt.cfg, cfg, o))
		}
	}
}
