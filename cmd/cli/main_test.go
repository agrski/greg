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
	type args struct {
		host string
		org  string
		repo string
		url  string
	}
	tests := []struct {
		name string
		args args
		want location
		err  bool
	}{
		{
			name: "fail if any required arg is missing",
			args: args{},
			want: location{},
			err:  true,
		},
		{
			name: "fail if using org without repo",
			args: args{org: "fakeOrg"},
			want: location{},
			err:  true,
		},
		{
			name: "fail if using repo without org",
			args: args{repo: "fakeRepo"},
			want: location{},
			err:  true,
		},
		{
			name: "fail if using url with org",
			args: args{url: "https://github.com/fakeOrg/fakeRepo", org: "fakeOrg"},
			want: location{},
			err:  true,
		},
		{
			name: "fail if using url with repo",
			args: args{url: "https://github.com/fakeOrg/fakeRepo", repo: "fakeRepo"},
			want: location{},
			err:  true,
		},
		{
			name: "org and repo both provided",
			args: args{host: githubHost, org: "fakeOrg", repo: "fakeRepo"},
			want: location{host: githubHost, organisation: "fakeOrg", repository: "fakeRepo"},
			err:  false,
		},
		{
			name: "url provided",
			args: args{url: "https://github.com/fakeOrg/fakeRepo"},
			want: location{host: githubHost, organisation: "fakeOrg", repository: "fakeRepo"},
			err:  false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				hostFlag = tt.args.host
				orgFlag = tt.args.org
				repoFlag = tt.args.repo
				urlFlag = tt.args.url

				got, err := getLocation()
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
			args: args{l: location{host: "github.com", organisation: "agrski", repository: "gitfind"}},
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
			want: location{host: "github.com", organisation: "agrski", repository: "gitfind"},
			err:  false,
		},
		{
			name: "full URL with git suffix",
			args: args{rawURL: "https://github.com/agrski/gitfind.git"},
			want: location{host: "github.com", organisation: "agrski", repository: "gitfind"},
			err:  false,
		},
		{
			name: "without scheme",
			args: args{rawURL: "github.com/agrski/gitfind"},
			want: location{host: "github.com", organisation: "agrski", repository: "gitfind"},
			err:  false,
		},
		{
			name: "with extra path",
			args: args{rawURL: "https://github.com/agrski/gitfind/tree/master/pkg"},
			want: location{host: "github.com", organisation: "agrski", repository: "gitfind"},
			err:  false,
		},
		{
			name: "full URL - GitLab",
			args: args{rawURL: "https://gitlab.com/agrski/gitfind"},
			want: location{host: "gitlab.com", organisation: "agrski", repository: "gitfind"},
			err:  false,
		},
		{
			name: "missing repo - trailing slash",
			args: args{rawURL: "https://github.com/agrski/"},
			want: location{},
			err:  true,
		},
		{
			name: "missing repo - no trailing slash",
			args: args{rawURL: "https://github.com/agrski"},
			want: location{},
			err:  true,
		},
		{
			name: "missing org and repo",
			args: args{rawURL: "https://github.com/"},
			want: location{},
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
