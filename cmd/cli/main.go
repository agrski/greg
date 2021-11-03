package main

import (
	"flag"
	"fmt"
	"log"
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

var (
	orgFlag string
	repoFlag string
)

func parseArguments() {
	flag.StringVar(&orgFlag, "org", "", "GitHub organisation, e.g. agrski")
	flag.StringVar(&repoFlag, "repo", "", "GitHub repository, e.g. gitfind")
	flag.Parse()
}

func getLocation() location {
	if isEmpty(orgFlag) {
		log.Fatal("org must be specified")
	}
	if isEmpty(repoFlag) {
		log.Fatal("repo must be specified")
	}

	return location{
		organisationName(orgFlag),
		repositoryName(repoFlag),
	}
}

func isEmpty(s string) bool {
	return "" == strings.TrimSpace(s)
}

func makeURI(l location) url.URL {
	return url.URL{
		Scheme: httpScheme,
		Host: githubHost,
		Path: fmt.Sprintf("%s/%s", l.organisation, l.repository),
	}
}

func main() {
	parseArguments()
	l := getLocation()
	u := makeURI(l)
	fmt.Printf("Retrieving files from %s", u.String())
}
