package main

import (
	"flag"
	"fmt"
	"net/url"
	"strings"
)

type organisationName string
type repositoryName string

type location struct {
	organisation organisationName
	repository   repositoryName
}

const (
	httpScheme = "https"
	githubHost = "github.com"
)

func getLocation() location {
	org := flag.String("org", "", "GitHub organisation, e.g. agrski")
	repo := flag.String("repo", "", "GitHub repository, e.g. gitfind")

	if isEmpty(org) {
		fmt.Errorf("org must be specified")
	}
	if isEmpty(repo) {
		fmt.Errorf("repo must be specified")
	}

	return location{
		organisationName(*org),
		repositoryName(*repo),
	}
}

func isEmpty(s *string) bool {
	return s == nil	|| strings.TrimSpace(*s) == ""
}

func makeURI(l location) url.URL {
	return url.URL{
		Scheme: httpScheme,
		Host: githubHost,
		Path: fmt.Sprintf("%s/%s", l.organisation, l.repository),
	}
}

func main() {
	l := getLocation()
	u := makeURI(l)
	fmt.Printf("Retrieving files from %s", u.String())
}
