// Package libdns-loopia implements a DNS record management client compatible
// with the libdns interfaces for Loopia.
//go:build integration

package loopia

import (
	"context"
	"net/netip"
	"reflect"
	"testing"
	"time"

	"github.com/libdns/libdns"
)

func getRecords() []libdns.Record {
	ip421 := netip.MustParseAddr("192.168.42.1")
	ip422 := netip.MustParseAddr("192.168.42.2")
	return []libdns.Record{
		libdns.Address{Name: "*", IP: ip421, TTL: time.Duration(5 * int(time.Minute))},
		libdns.Address{Name: "*", IP: ip422, TTL: time.Duration(5 * int(time.Minute))},
		libdns.NS{Name: "@", Target: "ns1.test.local.", TTL: time.Duration(int(time.Hour))},
		libdns.NS{Name: "@", Target: "ns2.test.local.", TTL: time.Duration(10 * int(time.Minute))},
		libdns.Address{Name: "www", IP: netip.MustParseAddr("1.1.1.1"), TTL: time.Duration(5 * int(time.Minute))},
		libdns.TXT{Name: "_challenge.test", Text: "foo", TTL: 0},
	}
}

func TestProvider_GetRecords(t *testing.T) {
	tc := setupTest(t)
	defer teardownTest(tc)

	type args struct {
		ctx  context.Context
		zone string
	}
	tests := []struct {
		name     string
		provider *Provider
		args     args
		want     []libdns.Record
		wantErr  bool
	}{
		{"first", tc.getProvider(), args{context.TODO(), "test.local"}, getRecords(), false},
		{"subdomain", tc.getProvider(), args{context.TODO(), "test.test.local"}, []libdns.Record{
			libdns.TXT{Name: "_challenge", Text: "foo", TTL: 0},
		}, false},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := tt.provider
			got, err := p.GetRecords(tt.args.ctx, tt.args.zone)
			if (err != nil) != tt.wantErr {
				t.Errorf("Provider.GetRecords() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Provider.GetRecords()\n got\t %v,\n want\t %v", got, tt.want)
			}
		})
	}
}

func TestProvider_AppendRecords(t *testing.T) {
	tc := setupTest(t)
	defer teardownTest(tc)
	type args struct {
		ctx     context.Context
		zone    string
		records []libdns.Record
	}
	tests := []struct {
		name     string
		provider *Provider
		args     args
		want     []libdns.Record
		wantErr  bool
	}{
		{"cdn", tc.getProvider(), args{context.TODO(), "test.local", []libdns.Record{
			libdns.TXT{Name: "_test", Text: "some text", TTL: time.Duration(5 * time.Minute)},
		}}, []libdns.Record{libdns.TXT{Name: "_test", Text: "some text", TTL: time.Duration(5 * time.Minute)}}, false},
		{"acme", tc.getProvider(),
			args{
				context.TODO(),
				"test.test.local",
				[]libdns.Record{
					libdns.TXT{Name: "_challenge", Text: "foo"},
				},
			},
			[]libdns.Record{libdns.TXT{Name: "_challenge", Text: "foo", TTL: 0}},
			false,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := tt.provider
			got, err := p.AppendRecords(tt.args.ctx, tt.args.zone, tt.args.records)
			if (err != nil) != tt.wantErr {
				t.Errorf("Provider.AppendRecords() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Provider.AppendRecords() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProvider_SetRecords(t *testing.T) {
	tc := setupTest(t)
	defer teardownTest(tc)

	type args struct {
		ctx     context.Context
		zone    string
		records []libdns.Record
	}
	tests := []struct {
		name     string
		provider *Provider
		args     args
		want     []libdns.Record
		wantErr  bool
	}{
		{"nil records", tc.getProvider(), args{context.TODO(), "test.local", nil}, nil, true},
		{"empty records", tc.getProvider(), args{context.TODO(), "test.local", []libdns.Record{}}, nil, true},
		{"invalid record", tc.getProvider(), args{context.TODO(), "test.local", []libdns.Record{libdns.Address{Name: "www"}}}, nil, true},
		{"invalid ID", tc.getProvider(), args{context.TODO(), "test.local", []libdns.Record{libdns.Address{Name: "www", IP: netip.MustParseAddr("127.0.0.1"), TTL: 5 * time.Minute}}}, nil, true},
		{"valid record", tc.getProvider(), args{context.TODO(), "test.local", []libdns.Record{libdns.Address{Name: "www", IP: netip.MustParseAddr("127.0.0.1"), TTL: 5 * time.Minute}}},
			[]libdns.Record{libdns.Address{Name: "www", IP: netip.MustParseAddr("127.0.0.1"), TTL: 5 * time.Minute}}, false},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := tt.provider
			got, err := p.SetRecords(tt.args.ctx, tt.args.zone, tt.args.records)
			if (err != nil) != tt.wantErr {
				t.Errorf("Provider.SetRecords() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Provider.SetRecords() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProvider_DeleteRecords(t *testing.T) {
	tc := setupTest(t)
	defer teardownTest(tc)

	type args struct {
		ctx     context.Context
		zone    string
		records []libdns.Record
	}
	tests := []struct {
		name     string
		provider *Provider
		args     args
		want     []libdns.Record
		wantErr  bool
	}{
		{"invalid zone", tc.getProvider(), args{context.TODO(), "", nil}, nil, true},
		{"nil records", tc.getProvider(), args{context.TODO(), "test.local", nil}, nil, true},
		{"empty records", tc.getProvider(), args{context.TODO(), "test.local", []libdns.Record{}}, nil, true},
		{"no id records", tc.getProvider(), args{context.TODO(), "test.local", []libdns.Record{libdns.Address{Name: "test"}}}, nil, true},
		// {"valid records", tc.getProvider(), args{context.TODO(), "test.local", []libdns.Record{{Name: "test", ID: "12345"}}}, []libdns.Record{{Name: "test", ID: "12345"}}, false},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := tt.provider
			got, err := p.DeleteRecords(tt.args.ctx, tt.args.zone, tt.args.records)
			if (err != nil) != tt.wantErr {
				t.Errorf("Provider.DeleteRecords() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Provider.DeleteRecords() = %v, want %v", got, tt.want)
			}
		})
	}
}
