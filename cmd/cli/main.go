package main

import (
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

func getLocation() location {
	if isEmpty(urlFlag) && (isEmpty(orgFlag) || isEmpty(repoFlag)) {
		log.Fatal("must specify either url or both org and repo")
	}

	if isEmpty(urlFlag) {
		return location{
			organisationName(orgFlag),
			repositoryName(repoFlag),
		}
	}

	return parseLocationFromUrl(urlFlag)
}

func parseLocationFromUrl(rawUrl string) location {
	u, err := url.Parse(rawUrl)
	if err != nil {
		log.Fatal("unable to parse URL", err)
	}

	p := u.Path
	p = strings.TrimPrefix(p, "/")
	orgAndRepo := strings.SplitN(p, "/", 3)
	if len(orgAndRepo) < 2 || isEmpty(orgAndRepo[0]) || isEmpty(orgAndRepo[1]) {
		log.Fatalf("unable to extract both org and repo from %s", u)
	}

	return location{
		organisationName(orgAndRepo[0]),
		repositoryName(orgAndRepo[1]),
	}
}

func getFiletypes() []string {
	suffixes := strings.Split(filetypeFlag, ",")
	withoutDots := make([]string, len(suffixes))
	for _, s := range suffixes {
		withoutDots = append(withoutDots, strings.TrimPrefix(s, "."))
	}
	return withoutDots
}

func getSearchPattern() string {
	if isEmpty(pattern) {
		log.Fatal("search term must be specified; wrap multiple words in quotes")
	}
	return pattern
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
	l := getLocation()
	u := makeURI(l)
	p := getSearchPattern()

	fmt.Printf("Searching for %s in %s", p, u.String())
}
