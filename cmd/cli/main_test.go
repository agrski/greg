//go:build !integration

package main

import (
	"net/url"
	"reflect"
	"testing"

	"github.com/agrski/gitfind/pkg/fetch"
)

func Test_makeURI(t *testing.T) {
	type test struct {
		name     string
		location fetch.Location
		want     url.URL
	}

	tests := []test{
		{
			name:     "github.com/agrski/gitfind",
			location: fetch.Location{Host: "github.com", Organisation: "agrski", Repository: "gitfind"},
			want: url.URL{
				Scheme: "https",
				Host:   "github.com",
				Path:   "agrski/gitfind",
				User:   nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := makeURI(tt.location); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("makeURI() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
