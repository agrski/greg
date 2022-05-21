//go:build !integration

package main

import (
	"net/url"
	"testing"

	"github.com/agrski/gitfind/pkg/fetch"
	"github.com/stretchr/testify/require"
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
				actual := makeURI(tt.location)
				require.Equal(t, tt.want, actual)
			},
		)
	}
}
