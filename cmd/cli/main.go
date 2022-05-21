package main

import (
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/agrski/greg/pkg/fetch"
	"github.com/rs/zerolog"
)

func main() {
	logger := makeLogger(zerolog.InfoLevel)

	args, err := GetArgs()
	if err != nil {
		logger.Fatal().Err(err).Send()
	}

	fetcher := fetch.New(logger, args.location, args.tokenSource)
	uri := makeURI(args.location)

	logger.
		Info().
		Str("pattern", args.searchPattern).
		Str("URL", uri.String()).
		Msg("searching")

	fetcher.Start()
	next, ok := fetcher.Next()
	if ok {
		fmt.Println(next)
	}
	fetcher.Stop()
}

func makeLogger(level zerolog.Level) zerolog.Logger {
	fieldKeyFormatter := func(v interface{}) string {
		return strings.ToUpper(
			fmt.Sprintf("%s=", v),
		)
	}
	logWriter := zerolog.ConsoleWriter{
		Out:        os.Stderr,
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

func makeURI(l fetch.Location) url.URL {
	return url.URL{
		Scheme: httpScheme,
		Host:   string(l.Host),
		Path:   fmt.Sprintf("%s/%s", l.Organisation, l.Repository),
	}
}
