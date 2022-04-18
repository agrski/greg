//go:build integration

package github

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/agrski/gitfind/pkg/auth"
	"golang.org/x/oauth2"
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
				QueryParams{
					RepoOwner: tt.owner,
					RepoName:  tt.repo,
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
