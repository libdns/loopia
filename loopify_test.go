package loopia

import (
	"testing"
)

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func Test_loopify(t *testing.T) {
	type args struct {
		name string
		zone string
	}
	tests := []struct {
		name     string
		args     args
		wantName string
		wantZone string
	}{
		{"simple", args{"some", "example.org"}, "some", "example.org"},
		{"ending-dot", args{"some", "example.org."}, "some", "example.org."},
		{"complex-left", args{"some.lcl", "example.org"}, "some.lcl", "example.org"},
		{"complex-right", args{"some", "lcl.example.org"}, "some.lcl", "example.org"},
		{"complex-right-dot", args{"some", "lcl.example.org."}, "some.lcl", "example.org."},
		{"simple-blank-name", args{"", "example.org"}, "", "example.org"},
		{"complex-blank-name", args{"", "lcl.example.org"}, ".lcl", "example.org"},
		{"asdf", args{"", "stuff.lcl.example.org"}, ".stuff.lcl", "example.org"},
		{"asdf", args{"some", "stuff.lcl.example.org"}, "some.stuff.lcl", "example.org"},

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := loopify(tt.args.name, tt.args.zone)
			if got != tt.wantName {
				t.Errorf("loopify() gotName = %v, wantName %v", got, tt.wantName)
			}
			if got1 != tt.wantZone {
				t.Errorf("loopifyFQDN() gotZone = %v, wantZone %v", got1, tt.wantZone)
			}
		})
	}
}

func Test_unLoopify(t *testing.T) {
	type args struct {
		name string
		zone string
	}
	tests := []struct {
		name     string
		args     args
		nameWant string
		zoneWant string
	}{
		{"simple", args{"some", "example.org"}, "some", "example.org"},
		{"ending-dot", args{"some", "example.org."}, "some", "example.org."},
		{"complex-left", args{"some", "lcl.example.org"}, "some", "lcl.example.org"},
		{"complex-right", args{"some", "lcl.example.org"}, "some", "lcl.example.org"},
		{"complex-right-dot", args{"some", "lcl.example.org."}, "some", "lcl.example.org."},
		{"a", args{"some.lcl", "example.org"}, "some", "lcl.example.org"},
		{"b", args{"some.lcl", "example.org."}, "some", "lcl.example.org."},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := unLoopify(tt.args.name, tt.args.zone)
			if got != tt.nameWant {
				t.Errorf("unLoopify() name got = %v, want %v", got, tt.nameWant)
			}
			if got1 != tt.zoneWant {
				t.Errorf("unLoopify() zone got = %v, want %v", got1, tt.zoneWant)
			}
		})
	}
}
