package main

import (
	"errors"
	"flag"
	"fmt"
	"net/url"
	"strings"

	"github.com/agrski/gitfind/pkg/auth"
	"github.com/agrski/gitfind/pkg/fetch"
	"golang.org/x/oauth2"
)

const (
	httpScheme = "https"
	githubHost = "github.com"
)

var supportedHosts = [...]fetch.HostName{githubHost}

var (
	hostFlag        string
	orgFlag         string
	repoFlag        string
	urlFlag         string
	filetypeFlag    string
	pattern         string
	accessToken     string
	accessTokenFile string
)

type rawArgs struct {
	host            string
	org             string
	repo            string
	url             string
	filetype        string
	searchPattern   string
	accessToken     string
	accessTokenFile string
}

type Args struct {
	location        fetch.Location
	searchPattern   string
	filetypes       []string
	tokenSource     oauth2.TokenSource
	accessToken     string
	accessTokenFile string
}

func GetArgs() (*Args, error) {
	raw := parseArguments()

	allowed := isSupportedHost(fetch.HostName(raw.host))
	if !allowed {
		return nil, fmt.Errorf("unsupported git hosting provider %s", raw.host)
	}

	location, err := getLocation(raw)
	if err != nil {
		return nil, err
	}

	pattern, err := getSearchPattern(raw)
	if err != nil {
		return nil, err
	}

	tokenSource, err := getAccessToken(raw.accessToken, raw.accessTokenFile)
	if err != nil {
		return nil, err
	}

	filetypes := getFiletypes(raw)

	return &Args{
		location:      location,
		searchPattern: pattern,
		filetypes:     filetypes,
		tokenSource:   tokenSource,
	}, nil
}

func parseArguments() *rawArgs {
	args := rawArgs{}

	flag.StringVar(&args.host, "host", githubHost, "git hostname, default: github.com")
	flag.StringVar(&args.org, "org", "", "organisation name, e.g. agrski")
	flag.StringVar(&args.repo, "repo", "", "repository name, e.g. gitfind")
	flag.StringVar(
		&args.url,
		"url",
		"",
		"Full URL of git repository, e.g https://github.com/agrski/gitfind",
	)
	flag.StringVar(&args.filetype, "type", "", "filetype suffix, e.g. md or go")
	flag.StringVar(&args.accessToken, "access-token", "", "raw access token for repository access")
	flag.StringVar(
		&args.accessTokenFile,
		"access-token-file",
		"",
		"file containing access token for repository access",
	)
	flag.Parse()

	if 1 == flag.NArg() {
		args.searchPattern = flag.Arg(0)
	}

	return &args
}

func getLocation(args *rawArgs) (fetch.Location, error) {
	if isEmpty(args.url) && (isEmpty(args.org) || isEmpty(args.repo)) {
		return fetch.Location{}, errors.New("must specify either url or both org and repo")
	}

	if !isEmpty(args.url) && (!isEmpty(args.org) || !isEmpty(args.repo)) {
		return fetch.Location{}, errors.New("cannot specify both url and org or repo")
	}

	if isEmpty(args.url) {
		return fetch.Location{
			Host:         fetch.HostName(args.host),
			Organisation: fetch.OrganisationName(args.org),
			Repository:   fetch.RepositoryName(args.repo),
		}, nil
	}

	return parseLocationFromURL(args.url)
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

func getFiletypes(args *rawArgs) []string {
	if isEmpty(args.filetype) {
		return nil
	}

	suffixes := strings.Split(args.filetype, ",")
	for idx, s := range suffixes {
		withoutWhitespace := strings.TrimSpace(s)
		withoutLeadingDot := strings.TrimPrefix(withoutWhitespace, ".")

		suffixes[idx] = withoutLeadingDot
	}

	return suffixes
}

func getSearchPattern(args *rawArgs) (string, error) {
	if isEmpty(args.searchPattern) {
		return "", errors.New("search term must be specified; wrap multiple words in quotes")
	}
	return args.searchPattern, nil
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

func getAccessToken(rawAccessToken string, accessTokenFile string) (oauth2.TokenSource, error) {
	if isEmpty(accessToken) && isEmpty(accessTokenFile) {
		return nil, errors.New("must specify either access token or access token file")
	}

	if !isEmpty(accessToken) && !isEmpty(accessTokenFile) {
		return nil, errors.New("only one of access token and access token file may be specified")
	}

	tokenSource, err := auth.TokenSourceFromString(rawAccessToken)
	if err != nil {
		tokenSource, err = auth.TokenSourceFromFile(accessTokenFile)
	}

	return tokenSource, err
}
