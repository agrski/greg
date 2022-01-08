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
	type args struct {
		p string
	}
	tests := []struct {
		name string
		args args
		want string
		err  bool
	}{
		{
			name: "should return pattern when provided",
			args: args{p: "hello"},
			want: "hello",
			err:  false,
		},
		{
			name: "should fail when pattern is not provided",
			args: args{p: ""},
			want: "",
			err:  true,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				pattern = tt.args.p

				got, err := getSearchPattern()
				if tt.err != (err != nil) {
					t.Errorf("getSearchPattern() = \"'%v', %v\", want \"'%v', %v\"", got, err, tt.want, tt.err)
				}
				if got != tt.want {
					t.Errorf("getSearchPattern() = \"'%v', %v\", want \"'%v', %v\"", got, err, tt.want, tt.err)
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
		{
			name: "empty-string should be empty",
			args: args{s: ""},
			want: true,
		},
		{
			name: "spaces should be empty",
			args: args{s: "   "},
			want: true,
		},
		{
			name: "tabs should be empty",
			args: args{s: "\t\t"},
			want: true,
		},
		{
			name: "mixed whitespace should be empty",
			args: args{s: " \t \u00a0 \u2003"},
			want: true,
		},
		{
			name: "word should not be empty",
			args: args{s: "hello"},
			want: false,
		},
		{
			name: "multiple words should not be empty",
			args: args{s: "hello world"},
			want: false,
		},
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
		{
			name: "github.com/agrski/gitfind",
			args: args{l: location{organisation: "agrski", repository: "gitfind"}},
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
		err  bool
	}{
		{
			name: "full URL",
			args: args{rawURL: "https://github.com/agrski/gitfind"},
			want: location{organisation: "agrski", repository: "gitfind"},
			err:  false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got, err := parseLocationFromURL(tt.args.rawURL)
				gotErr := err != nil

				if tt.err != gotErr {
					t.Errorf("parseLocationFromURL() = %v %v, want %v %v", got, gotErr, tt.want, tt.err)
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("parseLocationFromURL() = %v %v, want %v %v", got, gotErr, tt.want, tt.err)
				}
			},
		)
	}
}
