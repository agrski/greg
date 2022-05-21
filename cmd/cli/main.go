package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/agrski/gitfind/pkg/auth"
	"github.com/agrski/gitfind/pkg/fetch"
	"github.com/rs/zerolog"
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
	flag.StringVar(&accessToken, "access-token", "", "raw access token for repository access")
	flag.StringVar(
		&accessTokenFile,
		"access-token-file",
		"",
		"file containing access token for repository access",
	)
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

func makeLogger(level zerolog.Level) *zerolog.Logger {
	fieldKeyFormatter := func(v interface{}) string {
		return strings.ToUpper(
			fmt.Sprintf("%s=", v),
		)
	}
	logWriter := zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.RFC3339,
		FormatLevel: func(v interface{}) string {
			return strings.ToUpper(
				fmt.Sprintf("%-6s ", v),
			)
		},
		FormatFieldName:    fieldKeyFormatter,
		FormatErrFieldName: fieldKeyFormatter,
	}
	logger := zerolog.
		New(logWriter).
		Level(level).
		With().
		Timestamp().
		Logger()

	return &logger
}

func main() {
	logger := makeLogger(zerolog.InfoLevel)
	log.SetOutput(os.Stderr)

	parseArguments()

	l, err := getLocation()
	if err != nil {
		logger.Fatal().Err(err).Send()
	}

	allowed := isSupportedHost(l.Host)
	if !allowed {
		logger.
			Fatal().
			Err(fmt.Errorf("unsupported git hosting provider %s", l.Host)).
			Send()
	}

	u := makeURI(l)
	p, err := getSearchPattern()
	if err != nil {
		logger.Fatal().Err(err).Send()
	}

	tokenSource, err := getAccessToken(accessToken, accessTokenFile)
	if err != nil {
		logger.Fatal().Err(err).Send()
	}

	fetcher := fetch.New(l, tokenSource)

	fetcher.Start()
	fmt.Printf("Searching for %s in %s\n", p, u.String())

	next, ok := fetcher.Next()
	if ok {
		fmt.Println(next)
	}

	fetcher.Stop()
}
