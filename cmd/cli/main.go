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

type hostName string
type organisationName string
type repositoryName string

type location struct {
	host         hostName
	organisation organisationName
	repository   repositoryName
}

const (
	httpScheme = "https"
	githubHost = "github.com"
)

var supportedHosts = [...]hostName{githubHost}

var (
	hostFlag     string
	orgFlag      string
	repoFlag     string
	urlFlag      string
	filetypeFlag string
	pattern      string
)

func parseArguments() {
	flag.StringVar(&hostFlag, "host", githubHost, "Git-hosting hostname, default: github.com")
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

func getLocation() (location, error) {
	if isEmpty(urlFlag) && (isEmpty(orgFlag) || isEmpty(repoFlag)) {
		return location{}, errors.New("must specify either url or both org and repo")
	}

	if !isEmpty(urlFlag) && (!isEmpty(orgFlag) || !isEmpty(repoFlag)) {
		return location{}, errors.New("cannot specify both url and org or repo")
	}

	if isEmpty(urlFlag) {
		return location{
			hostName(hostFlag),
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

	if isEmpty(host) {
		return location{}, errors.New("host cannot be empty")
	}
	if isEmpty(org) {
		return location{}, errors.New("org cannot be empty")
	}
	if isEmpty(repo) {
		return location{}, errors.New("repo cannot be empty")
	}

	return location{
		hostName(host),
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
		Host:   string(l.host),
		Path:   fmt.Sprintf("%s/%s", l.organisation, l.repository),
	}
}

func isSupportedHost(host hostName) bool {
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

	allowed := isSupportedHost(l.host)
	if !allowed {
		log.Fatalf("unsupported git hosting provider %s", l.host)
	}

	u := makeURI(l)
	p, err := getSearchPattern()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Searching for %s in %s", p, u.String())
}
