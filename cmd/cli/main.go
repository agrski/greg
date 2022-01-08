package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/url"
	"strings"
)

// TODO
// 	- add tests for functionality here

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
	orgFlag      string
	repoFlag     string
	urlFlag      string
	filetypeFlag string
	pattern      string
)

func parseArguments() {
	flag.StringVar(&orgFlag, "org", "", "GitHub organisation, e.g. agrski")
	flag.StringVar(&repoFlag, "repo", "", "GitHub repository, e.g. gitfind")
	flag.StringVar(
		&urlFlag,
		"url",
		"",
		"Full URL of GitHub repository, e.g https://github.com/agrski/gitfind",
	)
	flag.StringVar(&filetypeFlag, "type", "", "filetype suffix, e.g. md or go")
	flag.Parse()
	if 1 == flag.NArg() {
		pattern = flag.Arg(0)
	}
}

func getLocation() (location, error) {
	if isEmpty(urlFlag) && (isEmpty(orgFlag) || isEmpty(repoFlag)) {
		return location{}, errors.New("must specify either url or both org and repo")
	}

	if isEmpty(urlFlag) {
		return location{
			organisationName(orgFlag),
			repositoryName(repoFlag),
		}, nil
	}

	return parseLocationFromURL(urlFlag)
}

func parseLocationFromURL(rawURL string) (location, error) {
	if isEmpty(rawURL) {
		return location{}, errors.New("cannot parse empty string")
	}

	noWhitespace := strings.TrimSpace(rawURL)

	parts := strings.SplitAfter(noWhitespace, "://")
	if len(parts) > 2 {
		return location{}, fmt.Errorf("cannot parse malformed string '%v'", noWhitespace)
	}

	withoutScheme := parts[len(parts)-1]
	hostAndPath := strings.Split(withoutScheme, "/")
	if len(hostAndPath) < 3 {
		return location{}, fmt.Errorf("unable to parse host, org, and repo from %v", hostAndPath)
	}

	host := hostAndPath[0]
	org := hostAndPath[1]
	repo := hostAndPath[2]
	repo = strings.TrimSuffix(repo, ".git")

	if !strings.HasSuffix(host, "github.com") {
		return location{}, fmt.Errorf("require github.com not %v", host)
	}
	if isEmpty(org) {
		return location{}, errors.New("org cannot be empty")
	}
	if isEmpty(repo) {
		return location{}, errors.New("repo cannot be empty")
	}

	return location{
		organisationName(org),
		repositoryName(repo),
	}, nil
}

func getFiletypes() []string {
	if isEmpty(filetypeFlag) {
		return nil
	}

	suffixes := strings.Split(filetypeFlag, ",")
	for idx, s := range suffixes {
		withoutWhitespace := strings.TrimSpace(s)
		withoutLeadingDot := strings.TrimPrefix(withoutWhitespace, ".")

		suffixes[idx] = withoutLeadingDot
	}

	return suffixes
}

func getSearchPattern() (string, error) {
	if isEmpty(pattern) {
		return "", errors.New("search term must be specified; wrap multiple words in quotes")
	}
	return pattern, nil
}

func isEmpty(s string) bool {
	return "" == strings.TrimSpace(s)
}

func makeURI(l location) url.URL {
	return url.URL{
		Scheme: httpScheme,
		Host:   githubHost,
		Path:   fmt.Sprintf("%s/%s", l.organisation, l.repository),
	}
}

func main() {
	parseArguments()
	l, err := getLocation()
	if err != nil {
		log.Fatal(err)
	}
	u := makeURI(l)
	p, err := getSearchPattern()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Searching for %s in %s", p, u.String())
}
