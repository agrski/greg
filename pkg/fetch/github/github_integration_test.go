//go:build integration

package github

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"

	"github.com/agrski/greg/pkg/auth"
	"github.com/agrski/greg/pkg/fetch/types"
	common "github.com/agrski/greg/pkg/types"
)

func getTokenSource(t *testing.T) oauth2.TokenSource {
	tokenPath := filepath.Join("testdata", "token.txt")
	source, err := auth.TokenSourceFromFile(tokenPath)
	if os.IsNotExist(err) {
		t.Errorf("please add a valid GitHub access token to %s", tokenPath)
	}

	return source
}

func TestGetDefaultBranchRef(t *testing.T) {
	type test struct {
		name     string
		owner    string
		repo     string
		expected string
	}

	tests := []test{
		{name: "agrski/gitfind", owner: "agrski", repo: "gitfind", expected: "master"},
		{name: "prometheus/prometheus", owner: "prometheus", repo: "prometheus", expected: "main"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := New(
				zerolog.Nop(),
				types.Location{
					Organisation: types.OrganisationName(tt.owner),
					Repository:   types.RepositoryName(tt.repo),
				},
				getTokenSource(t),
			)

			name, err := g.getDefaultBranchRef()
			if err != nil {
				t.Errorf("failed to query GitHub: %v", err)
			}

			if tt.expected != name {
				t.Errorf("expected %s but got %s", tt.expected, name)
			}
		})
	}
}

func TestEnsureCommitish(t *testing.T) {
	type test struct {
		name     string
		commit   string
		expected string
	}

	tests := []test{
		{name: "not provided so use default", commit: "", expected: "master"},
		{name: "should use provided", commit: "someHash", expected: "someHash"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := New(
				zerolog.Nop(),
				types.Location{
					Organisation: types.OrganisationName("agrski"),
					Repository:   types.RepositoryName("gitfind"),
				},
				getTokenSource(t),
			)

			g.queryParams.Commitish = tt.commit

			err := g.ensureCommitish()
			if err != nil {
				t.Error(err)
			}

			if tt.expected != g.queryParams.Commitish {
				t.Errorf("expected %s but got %s", tt.expected, g.queryParams.Commitish)
			}
		})
	}
}

func TestGetFiles(t *testing.T) {
	numResults := 0
	g := New(
		zerolog.Nop(),
		types.Location{
			Organisation: types.OrganisationName("agrski"),
			Repository:   types.RepositoryName("gitfind"),
		},
		getTokenSource(t),
	)

	fs, cancel := g.getFiles()
	for f := range fs {
		require.NotNil(t, f)
		numResults++
	}
	cancel()

	require.NotZero(t, numResults)
}

func TestStart(t *testing.T) {
	numFiles := 5
	results := make([]*common.FileInfo, 0, numFiles)
	g := New(
		zerolog.Nop(),
		types.Location{
			Organisation: types.OrganisationName("agrski"),
			Repository:   types.RepositoryName("gitfind"),
		},
		getTokenSource(t),
	)

	// Ensure API returns some files
	g.Start()
	for i := 0; i < numFiles; i++ {
		if next, ok := g.Next(); ok {
			results = append(results, next)
		}
	}
	g.Stop()

	require.Len(t, results, numFiles)
}

func TestStop(t *testing.T) {
	g := New(
		zerolog.Nop(),
		types.Location{
			Organisation: types.OrganisationName("agrski"),
			Repository:   types.RepositoryName("gitfind"),
		},
		getTokenSource(t),
	)

	// Stopping immediately should be far too fast for any real results to be fetched
	g.Start()
	g.Stop()

	next, ok := g.Next()
	require.False(t, ok)
	require.Nil(t, next)
	assert.Empty(t, g.results)
}
