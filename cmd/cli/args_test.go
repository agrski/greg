//go:build !integration

package main

import (
	"testing"

	"github.com/agrski/gitfind/pkg/fetch"
	"github.com/stretchr/testify/require"
)

func Test_getFiletypes(t *testing.T) {
	type test struct {
		name      string
		filetypes string
		want      []string
	}

	tests := []test{
		{
			name:      "no filetypes succeeds with nil",
			filetypes: "",
			want:      nil,
		},
		{
			name:      "simple letter-only filetype is unchanged",
			filetypes: "md",
			want:      []string{"md"},
		},
		{
			name:      "dot prefix is removed",
			filetypes: ".md",
			want:      []string{"md"},
		},
		{
			name:      "multiple suffices are handled correctly",
			filetypes: ".md,go,.txt",
			want:      []string{"md", "go", "txt"},
		},
		{
			name:      "multiple suffices with whitespace are handled correctly",
			filetypes: ".md,go , .txt",
			want:      []string{"md", "go", "txt"},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				actual := getFiletypes(tt.filetypes)
				require.ElementsMatch(t, tt.want, actual)
			},
		)
	}
}

func Test_getLocation(t *testing.T) {
	type test struct {
		name    string
		args    *rawArgs
		want    fetch.Location
		wantErr bool
	}

	tests := []test{
		{
			name:    "fail if any required arg is missing",
			args:    &rawArgs{},
			want:    fetch.Location{},
			wantErr: true,
		},
		{
			name:    "fail if using org without repo",
			args:    &rawArgs{org: "fakeOrg"},
			want:    fetch.Location{},
			wantErr: true,
		},
		{
			name:    "fail if using repo without org",
			args:    &rawArgs{repo: "fakeRepo"},
			want:    fetch.Location{},
			wantErr: true,
		},
		{
			name:    "fail if using url with org",
			args:    &rawArgs{url: "https://github.com/fakeOrg/fakeRepo", org: "fakeOrg"},
			want:    fetch.Location{},
			wantErr: true,
		},
		{
			name:    "fail if using url with repo",
			args:    &rawArgs{url: "https://github.com/fakeOrg/fakeRepo", repo: "fakeRepo"},
			want:    fetch.Location{},
			wantErr: true,
		},
		{
			name:    "org and repo both provided",
			args:    &rawArgs{host: githubHost, org: "fakeOrg", repo: "fakeRepo"},
			want:    fetch.Location{Host: githubHost, Organisation: "fakeOrg", Repository: "fakeRepo"},
			wantErr: false,
		},
		{
			name:    "url provided",
			args:    &rawArgs{url: "https://github.com/fakeOrg/fakeRepo"},
			want:    fetch.Location{Host: githubHost, Organisation: "fakeOrg", Repository: "fakeRepo"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				actual, err := getLocation(tt.args)

				if tt.wantErr {
					require.Error(t, err)
				} else {
					require.NoError(t, err)
				}
				require.Equal(t, tt.want, actual)
			},
		)
	}
}

func Test_getSearchPattern(t *testing.T) {
	type test struct {
		name          string
		searchPattern string
		want          string
		wantErr       bool
	}

	tests := []test{
		{
			name:          "should return pattern when provided",
			searchPattern: "hello",
			want:          "hello",
			wantErr:       false,
		},
		{
			name:          "should fail when pattern is not provided",
			searchPattern: "",
			want:          "",
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				actual, err := getSearchPattern(tt.searchPattern)

				if tt.wantErr {
					require.Error(t, err)
				} else {
					require.NoError(t, err)
				}
				require.Equal(t, tt.want, actual)
			},
		)
	}
}

func Test_isEmpty(t *testing.T) {
	type test struct {
		name  string
		input string
		want  bool
	}

	tests := []test{
		{
			name:  "empty-string should be empty",
			input: "",
			want:  true,
		},
		{
			name:  "spaces should be empty",
			input: "   ",
			want:  true,
		},
		{
			name:  "tabs should be empty",
			input: "\t\t",
			want:  true,
		},
		{
			name:  "mixed whitespace should be empty",
			input: " \t \u00a0 \u2003",
			want:  true,
		},
		{
			name:  "word should not be empty",
			input: "hello",
			want:  false,
		},
		{
			name:  "multiple words should not be empty",
			input: "hello world",
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				actual := isEmpty(tt.input)
				require.Equal(t, tt.want, actual)
			},
		)
	}
}

func Test_parseLocationFromURL(t *testing.T) {
	type test struct {
		name    string
		rawUrl  string
		want    fetch.Location
		wantErr bool
	}

	tests := []test{
		{
			name:    "full URL",
			rawUrl:  "https://github.com/agrski/gitfind",
			want:    fetch.Location{Host: "github.com", Organisation: "agrski", Repository: "gitfind"},
			wantErr: false,
		},
		{
			name:    "full URL with git suffix",
			rawUrl:  "https://github.com/agrski/gitfind.git",
			want:    fetch.Location{Host: "github.com", Organisation: "agrski", Repository: "gitfind"},
			wantErr: false,
		},
		{
			name:    "without scheme",
			rawUrl:  "github.com/agrski/gitfind",
			want:    fetch.Location{Host: "github.com", Organisation: "agrski", Repository: "gitfind"},
			wantErr: false,
		},
		{
			name:    "with extra path",
			rawUrl:  "https://github.com/agrski/gitfind/tree/master/pkg",
			want:    fetch.Location{Host: "github.com", Organisation: "agrski", Repository: "gitfind"},
			wantErr: false,
		},
		{
			name:    "full URL - GitLab",
			rawUrl:  "https://gitlab.com/agrski/gitfind",
			want:    fetch.Location{Host: "gitlab.com", Organisation: "agrski", Repository: "gitfind"},
			wantErr: false,
		},
		{
			name:    "missing repo - trailing slash",
			rawUrl:  "https://github.com/agrski/",
			want:    fetch.Location{},
			wantErr: true,
		},
		{
			name:    "missing repo - no trailing slash",
			rawUrl:  "https://github.com/agrski",
			want:    fetch.Location{},
			wantErr: true,
		},
		{
			name:    "missing org and repo",
			rawUrl:  "https://github.com/",
			want:    fetch.Location{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				actual, err := parseLocationFromURL(tt.rawUrl)

				if tt.wantErr {
					require.Error(t, err)
				} else {
					require.NoError(t, err)
				}
				require.Equal(t, tt.want, actual)
			},
		)
	}
}
