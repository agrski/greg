package auth

import (
	"golang.org/x/oauth2"
)

func ToTokenSource(rawToken string) oauth2.TokenSource {
	return oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: rawToken,
		},
	)
}
