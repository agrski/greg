//go:build !integration

package main

import (
	"net/url"
	"reflect"
	"testing"

	"github.com/agrski/gitfind/pkg/fetch"
)

func Test_getFiletypes(t *testing.T) {
	type test struct {
		name string
		args *rawArgs
		want []string
	}

	tests := []test{
		{
			name: "no filetypes succeeds with nil",
			args: &rawArgs{filetype: ""},
			want: nil,
		},
		{
			name: "simple letter-only filetype is unchanged",
			args: &rawArgs{filetype: "md"},
			want: []string{"md"},
		},
		{
			name: "dot prefix is removed",
			args: &rawArgs{filetype: ".md"},
			want: []string{"md"},
		},
		{
			name: "multiple suffices are handled correctly",
			args: &rawArgs{filetype: ".md,go,.txt"},
			want: []string{"md", "go", "txt"},
		},
		{
			name: "multiple suffices with whitespace are handled correctly",
			args: &rawArgs{filetype: ".md,go , .txt"},
			want: []string{"md", "go", "txt"},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := getFiletypes(tt.args); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("getFiletypes() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func Test_getLocation(t *testing.T) {
	type test struct {
		name string
		args *rawArgs
		want fetch.Location
		err  bool
	}
	tests := []test{
		{
			name: "fail if any required arg is missing",
			args: &rawArgs{},
			want: fetch.Location{},
			err:  true,
		},
		{
			name: "fail if using org without repo",
			args: &rawArgs{org: "fakeOrg"},
			want: fetch.Location{},
			err:  true,
		},
		{
			name: "fail if using repo without org",
			args: &rawArgs{repo: "fakeRepo"},
			want: fetch.Location{},
			err:  true,
		},
		{
			name: "fail if using url with org",
			args: &rawArgs{url: "https://github.com/fakeOrg/fakeRepo", org: "fakeOrg"},
			want: fetch.Location{},
			err:  true,
		},
		{
			name: "fail if using url with repo",
			args: &rawArgs{url: "https://github.com/fakeOrg/fakeRepo", repo: "fakeRepo"},
			want: fetch.Location{},
			err:  true,
		},
		{
			name: "org and repo both provided",
			args: &rawArgs{host: githubHost, org: "fakeOrg", repo: "fakeRepo"},
			want: fetch.Location{Host: githubHost, Organisation: "fakeOrg", Repository: "fakeRepo"},
			err:  false,
		},
		{
			name: "url provided",
			args: &rawArgs{url: "https://github.com/fakeOrg/fakeRepo"},
			want: fetch.Location{Host: githubHost, Organisation: "fakeOrg", Repository: "fakeRepo"},
			err:  false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got, err := getLocation(tt.args)
				gotErr := err != nil

				if tt.err != gotErr {
					t.Errorf("getLocation() = %v %v, want %v %v", got, err, tt.want, tt.err)
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("getLocation() = %v %v, want %v %v", got, err, tt.want, tt.err)
				}
			},
		)
	}
}

func Test_getSearchPattern(t *testing.T) {
	type test struct {
		name string
		args *rawArgs
		want string
		err  bool
	}
	tests := []test{
		{
			name: "should return pattern when provided",
			args: &rawArgs{searchPattern: "hello"},
			want: "hello",
			err:  false,
		},
		{
			name: "should fail when pattern is not provided",
			args: &rawArgs{searchPattern: ""},
			want: "",
			err:  true,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got, err := getSearchPattern(tt.args)
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
		l fetch.Location
	}
	tests := []struct {
		name string
		args args
		want url.URL
	}{
		{
			name: "github.com/agrski/gitfind",
			args: args{l: fetch.Location{Host: "github.com", Organisation: "agrski", Repository: "gitfind"}},
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
		want fetch.Location
		err  bool
	}{
		{
			name: "full URL",
			args: args{rawURL: "https://github.com/agrski/gitfind"},
			want: fetch.Location{Host: "github.com", Organisation: "agrski", Repository: "gitfind"},
			err:  false,
		},
		{
			name: "full URL with git suffix",
			args: args{rawURL: "https://github.com/agrski/gitfind.git"},
			want: fetch.Location{Host: "github.com", Organisation: "agrski", Repository: "gitfind"},
			err:  false,
		},
		{
			name: "without scheme",
			args: args{rawURL: "github.com/agrski/gitfind"},
			want: fetch.Location{Host: "github.com", Organisation: "agrski", Repository: "gitfind"},
			err:  false,
		},
		{
			name: "with extra path",
			args: args{rawURL: "https://github.com/agrski/gitfind/tree/master/pkg"},
			want: fetch.Location{Host: "github.com", Organisation: "agrski", Repository: "gitfind"},
			err:  false,
		},
		{
			name: "full URL - GitLab",
			args: args{rawURL: "https://gitlab.com/agrski/gitfind"},
			want: fetch.Location{Host: "gitlab.com", Organisation: "agrski", Repository: "gitfind"},
			err:  false,
		},
		{
			name: "missing repo - trailing slash",
			args: args{rawURL: "https://github.com/agrski/"},
			want: fetch.Location{},
			err:  true,
		},
		{
			name: "missing repo - no trailing slash",
			args: args{rawURL: "https://github.com/agrski"},
			want: fetch.Location{},
			err:  true,
		},
		{
			name: "missing org and repo",
			args: args{rawURL: "https://github.com/"},
			want: fetch.Location{},
			err:  true,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got, err := parseLocationFromURL(tt.args.rawURL)
				gotErr := err != nil

				if tt.err != gotErr {
					t.Errorf("parseLocationFromURL() = %v %v, want %v %v", got, err, tt.want, tt.err)
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("parseLocationFromURL() = %v %v, want %v %v", got, err, tt.want, tt.err)
				}
			},
		)
	}
}
