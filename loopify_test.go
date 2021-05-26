package loopia

import (
	"reflect"
	"testing"

	"github.com/libdns/libdns"
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

func Test_loopifyRecords(t *testing.T) {
	type args struct {
		zone    string
		records []libdns.Record
	}
	tests := []struct {
		name           string
		args           args
		wantHostSuffix string
		wantDomain     string
		wantOutNames   []string
	}{
		{"first", args{"lcl.test.local", getRecords()}, ".lcl", "test.local", []string{"*.lcl", "*.lcl", "@.lcl", "@.lcl", "www.lcl", "_challenge.test.lcl"}},
		// TODO: Add more test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHostSuffix, gotDomain := loopifyRecords(tt.args.zone, tt.args.records)
			if gotHostSuffix != tt.wantHostSuffix {
				t.Errorf("loopifyRecords() gotHostSuffix = %v, want %v", gotHostSuffix, tt.wantHostSuffix)
			}
			if gotDomain != tt.wantDomain {
				t.Errorf("loopifyRecords() gotDomain = %v, want %v", gotDomain, tt.wantDomain)
			}
			gotL := len(tt.args.records)
			wantL := len(tt.wantOutNames)
			min := minInt(gotL, wantL)
			for i := 0; i < min; i++ {
				if tt.args.records[i].Name != tt.wantOutNames[i] {
					t.Errorf("loopifyRecords got name = %v, want %v", tt.args.records[i].Name, tt.wantOutNames[i])
				}
			}
			if gotL != wantL {
				t.Errorf("loopifyRecords got = %v records, want %v", gotL, wantL)
			}
		})
	}
}

func Test_unLoopifyRecords(t *testing.T) {
	type args struct {
		hostSuffix string
		records    []libdns.Record
	}

	r1 := func() []libdns.Record {
		return []libdns.Record{
			{Name: "_challenge.test.lcl"},
		}
	}

	tests := []struct {
		name      string
		args      args
		wantNames []string
	}{
		{"first", args{".test", r1()}, []string{"_challenge.test.lcl"}},
		{"second", args{".lcl", r1()}, []string{"_challenge.test"}},
		{"third", args{".test.lcl", r1()}, []string{"_challenge"}},
		// TODO: Add more test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			unLoopifyRecords(tt.args.hostSuffix, tt.args.records)
			gotNames := make([]string, len(tt.args.records))
			for i, r := range tt.args.records {
				gotNames[i] = r.Name
			}
			if !reflect.DeepEqual(gotNames, tt.wantNames) {
				t.Errorf("unloopifyRecords() got names %v, want %v", gotNames, tt.wantNames)
			}

		})
	}
}
