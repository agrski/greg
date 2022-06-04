package main

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/agrski/greg/pkg/auth"
	fetchTypes "github.com/agrski/greg/pkg/fetch/types"
	"github.com/agrski/greg/pkg/match"
	"golang.org/x/oauth2"
)

type VerbosityLevel int

const (
	VerbosityQuiet VerbosityLevel = iota
	VerbosityNormal
	VerbosityHigh
)

const (
	httpScheme = "https"
	githubHost = "github.com"
)

var supportedHosts = [...]fetchTypes.HostName{githubHost}

type rawArgs struct {
	// Application behaviour
	host            string
	org             string
	repo            string
	url             string
	filetypes       string
	searchPattern   string
	accessToken     string
	accessTokenFile string
	// Presentation/display behaviour
	quiet   bool
	verbose bool
}

type Args struct {
	location      fetchTypes.Location
	searchPattern string
	filetypes     []string
	tokenSource   oauth2.TokenSource
	verbosity     VerbosityLevel
}

func GetArgs() (*Args, error) {
	raw, err := parseArguments()
	if err != nil {
		return nil, err
	}

	err = isSupportedHost(raw.host)
	if err != nil {
		return nil, err
	}

	location, err := getLocation(raw)
	if err != nil {
		return nil, err
	}

	pattern, err := getSearchPattern(raw.searchPattern)
	if err != nil {
		return nil, err
	}

	tokenSource, err := getAccessToken(raw.accessToken, raw.accessTokenFile)
	if err != nil {
		return nil, err
	}

	filetypes := getFiletypes(raw.filetypes)

	verbosity := getVerbosity(raw.quiet, raw.verbose)

	return &Args{
		location:      location,
		searchPattern: pattern,
		filetypes:     filetypes,
		tokenSource:   tokenSource,
		verbosity:     verbosity,
	}, nil
}

func parseArguments() (*rawArgs, error) {
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
	flag.StringVar(&args.filetypes, "type", "", "filetype suffix, e.g. md or go")
	flag.StringVar(&args.accessToken, "access-token", "", "raw access token for repository access")
	flag.StringVar(
		&args.accessTokenFile,
		"access-token-file",
		"",
		"file containing access token for repository access",
	)
	flag.BoolVar(&args.quiet, "quiet", false, "disable logging; overrides verbose mode")
	flag.BoolVar(&args.verbose, "verbose", false, "increase logging; overridden by quiet mode")
	flag.Parse()

	if 1 != flag.NArg() {
		return nil, fmt.Errorf("expected one search term but found %d", flag.NArg())
	}
	args.searchPattern = flag.Arg(0)

	return &args, nil
}

func getLocation(args *rawArgs) (fetchTypes.Location, error) {
	if isEmpty(args.url) && (isEmpty(args.org) || isEmpty(args.repo)) {
		return fetchTypes.Location{}, errors.New("must specify either url or both org and repo")
	}

	if !isEmpty(args.url) && (!isEmpty(args.org) || !isEmpty(args.repo)) {
		return fetchTypes.Location{}, errors.New("cannot specify both url and org or repo")
	}

	if isEmpty(args.url) {
		return fetchTypes.Location{
			Host:         fetchTypes.HostName(args.host),
			Organisation: fetchTypes.OrganisationName(args.org),
			Repository:   fetchTypes.RepositoryName(args.repo),
		}, nil
	}

	return parseLocationFromURL(args.url)
}

func parseLocationFromURL(rawURL string) (fetchTypes.Location, error) {
	if isEmpty(rawURL) {
		return fetchTypes.Location{}, errors.New("cannot parse empty string")
	}

	noWhitespace := strings.TrimSpace(rawURL)

	parts := strings.SplitAfter(noWhitespace, "://")
	if len(parts) > 2 {
		return fetchTypes.Location{}, fmt.Errorf("cannot parse malformed string '%v'", noWhitespace)
	}

	withoutScheme := parts[len(parts)-1]
	hostAndPath := strings.Split(withoutScheme, "/")
	if len(hostAndPath) < 3 {
		return fetchTypes.Location{}, fmt.Errorf("unable to parse host, org, and repo from %v", hostAndPath)
	}

	host := hostAndPath[0]
	org := hostAndPath[1]
	repo := hostAndPath[2]
	repo = strings.TrimSuffix(repo, ".git")

	if isEmpty(host) {
		return fetchTypes.Location{}, errors.New("host cannot be empty")
	}
	if isEmpty(org) {
		return fetchTypes.Location{}, errors.New("org cannot be empty")
	}
	if isEmpty(repo) {
		return fetchTypes.Location{}, errors.New("repo cannot be empty")
	}

	return fetchTypes.Location{
		Host:         fetchTypes.HostName(host),
		Organisation: fetchTypes.OrganisationName(org),
		Repository:   fetchTypes.RepositoryName(repo),
	}, nil
}

func getFiletypes(filetypes string) []string {
	if isEmpty(filetypes) {
		return nil
	}

	suffixes := strings.Split(filetypes, ",")
	for idx, s := range suffixes {
		normalised := match.NormaliseExtension(s)
		suffixes[idx] = normalised
	}

	return suffixes
}

func getSearchPattern(pattern string) (string, error) {
	if isEmpty(pattern) {
		return "", errors.New("search term must be specified; wrap multiple words in quotes")
	}
	return pattern, nil
}

func isEmpty(s string) bool {
	return "" == strings.TrimSpace(s)
}

func isSupportedHost(host string) error {
	hostname := fetchTypes.HostName(host)
	for _, h := range supportedHosts {
		if hostname == h {
			return nil
		}
	}
	return fmt.Errorf("unsupported git hosting provider %s", host)
}

func getAccessToken(
	accessToken string,
	accessTokenFile string,
) (oauth2.TokenSource, error) {
	if isEmpty(accessToken) && isEmpty(accessTokenFile) {
		return nil, errors.New("must specify either access token or access token file")
	}

	if !isEmpty(accessToken) && !isEmpty(accessTokenFile) {
		return nil, errors.New("only one of access token and access token file may be specified")
	}

	tokenSource, err := auth.TokenSourceFromString(accessToken)
	if err != nil {
		tokenSource, err = auth.TokenSourceFromFile(accessTokenFile)
	}

	return tokenSource, err
}

func getVerbosity(quiet bool, verbose bool) VerbosityLevel {
	if quiet {
		return VerbosityQuiet
	} else if verbose {
		return VerbosityHigh
	} else {
		return VerbosityNormal
	}
}
