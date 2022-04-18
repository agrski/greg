package auth

import (
	"errors"
	"os"
	"strings"

	"golang.org/x/oauth2"
)

func TokenSourceFromString(rawToken string) (oauth2.TokenSource, error) {
	if isEmpty(rawToken) {
		return nil, errors.New("access token cannot be empty")
	}

	return oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: rawToken,
		},
	), nil
}

func isEmpty(s string) bool {
	return strings.TrimSpace(s) == ""
}

func TokenSourceFromFile(fileName string) (oauth2.TokenSource, error) {
	if isEmpty(fileName) {
		return nil, errors.New("access token file name must be specified")
	}

	maybeToken, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	token := string(maybeToken)
	token = strings.TrimSpace(token)

	return TokenSourceFromString(token)
}
