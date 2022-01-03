package main

import (
	"net/url"
	"reflect"
	"testing"
)

func Test_getFiletypes(t *testing.T) {
	type args struct {
		filetypes string
	}

	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "no filetypes succeeds with nil",
			args: args{filetypes: ""},
			want: nil,
		},
		{
			name: "simple letter-only filetype is unchanged",
			args: args{filetypes: "md"},
			want: []string{"md"},
		},
		{
			name: "dot prefix is removed",
			args: args{filetypes: ".md"},
			want: []string{"md"},
		},
		{
			name: "multiple suffices are handled correctly",
			args: args{filetypes: ".md,go,.txt"},
			want: []string{"md", "go", "txt"},
		},
		{
			name: "multiple suffices with whitespace are handled correctly",
			args: args{filetypes: ".md,go , .txt"},
			want: []string{"md", "go", "txt"},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				filetypeFlag = tt.args.filetypes

				if got := getFiletypes(); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("getFiletypes() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func Test_getLocation(t *testing.T) {
	tests := []struct {
		name string
		want location
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := getLocation(); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("getLocation() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func Test_getSearchPattern(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := getSearchPattern(); got != tt.want {
					t.Errorf("getSearchPattern() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func Test_isEmpty(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := isEmpty(tt.args.s); got != tt.want {
					t.Errorf("isEmpty() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func Test_makeURI(t *testing.T) {
	type args struct {
		l location
	}
	tests := []struct {
		name string
		args args
		want url.URL
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := makeURI(tt.args.l); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("makeURI() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func Test_parseLocationFromURL(t *testing.T) {
	type args struct {
		rawURL string
	}
	tests := []struct {
		name string
		args args
		want location
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := parseLocationFromURL(tt.args.rawURL); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("parseLocationFromURL() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
