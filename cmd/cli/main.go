package main

import (
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"

	"github.com/agrski/greg/pkg/fetch"
	fetchTypes "github.com/agrski/greg/pkg/fetch/types"
	"github.com/agrski/greg/pkg/match"
	"github.com/agrski/greg/pkg/present/console"
)

func main() {
	args, err := GetArgs()
	if err != nil {
		logger := makeLogger(zerolog.InfoLevel, false)
		logger.Fatal().Err(err).Send()
	}

	logger := makeLogger(zerolog.InfoLevel, args.enableColour)

	switch args.verbosity {
	case VerbosityQuiet:
		logger = logger.Level(zerolog.Disabled)
	case VerbosityHigh:
		logger = logger.Level(zerolog.DebugLevel)
	case VerbosityNormal:
		// Already at normal (info-level) verbosity
	}

	console := console.New(os.Stdout, args.enableColour)

	matcher := match.New(logger, args.caseInsensitive, args.filetypes)

	fetcher := fetch.New(logger, args.location, args.tokenSource)
	uri := makeURI(args.location)

	logger.
		Info().
		Str("pattern", args.searchPattern).
		Str("URL", uri.String()).
		Msg("searching")

	err = fetcher.Start()
	if err != nil {
		logger.Fatal().Err(err).Msg("unable to start fetching matches")
	}

	next, ok := fetcher.Next()
	if ok {
		if m, ok := matcher.Match(args.searchPattern, next); ok {
			console.Write(next, m)
		}
	}

	_ = fetcher.Stop()
}

func makeLogger(level zerolog.Level, enableColour bool) zerolog.Logger {
	fieldKeyFormatter := func(v interface{}) string {
		return strings.ToUpper(
			fmt.Sprintf("%s=", v),
		)
	}
	logWriter := zerolog.ConsoleWriter{
		Out:        os.Stderr,
		NoColor:    !enableColour,
		TimeFormat: time.RFC3339,
		FormatLevel: func(v interface{}) string {
			l, ok := v.(string)
			switch {
			case !ok, l == zerolog.InfoLevel.String():
				return ""
			default:
				return strings.ToUpper(
					fmt.Sprintf("%-6s ", v),
				)
			}
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

	return logger
}

func makeURI(l fetchTypes.Location) url.URL {
	return url.URL{
		Scheme: httpScheme,
		Host:   string(l.Host),
		Path:   fmt.Sprintf("%s/%s", l.Organisation, l.Repository),
	}
}
