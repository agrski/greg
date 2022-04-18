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
