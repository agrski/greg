package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/agrski/gitfind/pkg/fetch"
)

const (
	httpScheme = "https"
	githubHost = "github.com"
)

var supportedHosts = [...]fetch.HostName{githubHost}

var (
	hostFlag     string
	orgFlag      string
	repoFlag     string
	urlFlag      string
	filetypeFlag string
	pattern      string
)

func parseArguments() {
	flag.StringVar(&hostFlag, "host", githubHost, "git hostname, default: github.com")
	flag.StringVar(&orgFlag, "org", "", "organisation name, e.g. agrski")
	flag.StringVar(&repoFlag, "repo", "", "repository name, e.g. gitfind")
	flag.StringVar(
		&urlFlag,
		"url",
		"",
		"Full URL of git repository, e.g https://github.com/agrski/gitfind",
	)
	flag.StringVar(&filetypeFlag, "type", "", "filetype suffix, e.g. md or go")
	flag.Parse()
	if 1 == flag.NArg() {
		pattern = flag.Arg(0)
	}
}

func getLocation() (fetch.Location, error) {
	if isEmpty(urlFlag) && (isEmpty(orgFlag) || isEmpty(repoFlag)) {
		return fetch.Location{}, errors.New("must specify either url or both org and repo")
	}

	if !isEmpty(urlFlag) && (!isEmpty(orgFlag) || !isEmpty(repoFlag)) {
		return fetch.Location{}, errors.New("cannot specify both url and org or repo")
	}

	if isEmpty(urlFlag) {
		return fetch.Location{
			Host:         fetch.HostName(hostFlag),
			Organisation: fetch.OrganisationName(orgFlag),
			Repository:   fetch.RepositoryName(repoFlag),
		}, nil
	}

	return parseLocationFromURL(urlFlag)
}

func parseLocationFromURL(rawURL string) (fetch.Location, error) {
	if isEmpty(rawURL) {
		return fetch.Location{}, errors.New("cannot parse empty string")
	}

	noWhitespace := strings.TrimSpace(rawURL)

	parts := strings.SplitAfter(noWhitespace, "://")
	if len(parts) > 2 {
		return fetch.Location{}, fmt.Errorf("cannot parse malformed string '%v'", noWhitespace)
	}

	withoutScheme := parts[len(parts)-1]
	hostAndPath := strings.Split(withoutScheme, "/")
	if len(hostAndPath) < 3 {
		return fetch.Location{}, fmt.Errorf("unable to parse host, org, and repo from %v", hostAndPath)
	}

	host := hostAndPath[0]
	org := hostAndPath[1]
	repo := hostAndPath[2]
	repo = strings.TrimSuffix(repo, ".git")

	if isEmpty(host) {
		return fetch.Location{}, errors.New("host cannot be empty")
	}
	if isEmpty(org) {
		return fetch.Location{}, errors.New("org cannot be empty")
	}
	if isEmpty(repo) {
		return fetch.Location{}, errors.New("repo cannot be empty")
	}

	return fetch.Location{
		Host:         fetch.HostName(host),
		Organisation: fetch.OrganisationName(org),
		Repository:   fetch.RepositoryName(repo),
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

func makeURI(l fetch.Location) url.URL {
	return url.URL{
		Scheme: httpScheme,
		Host:   string(l.Host),
		Path:   fmt.Sprintf("%s/%s", l.Organisation, l.Repository),
	}
}

func isSupportedHost(host fetch.HostName) bool {
	for _, h := range supportedHosts {
		if host == h {
			return true
		}
	}
	return false
}

func main() {
	parseArguments()

	l, err := getLocation()
	if err != nil {
		log.Fatal(err)
	}

	allowed := isSupportedHost(l.Host)
	if !allowed {
		log.Fatalf("unsupported git hosting provider %s", l.Host)
	}

	u := makeURI(l)
	p, err := getSearchPattern()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Searching for %s in %s", p, u.String())
}
