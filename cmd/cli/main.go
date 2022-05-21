package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/agrski/gitfind/pkg/fetch"
	"github.com/rs/zerolog"
)

func main() {
	logger := makeLogger(zerolog.InfoLevel)

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

	fetcher := fetch.New(logger, l, tokenSource)

	fetcher.Start()
	logger.Info().Str("pattern", p).Str("URL", u.String()).Msg("searching")

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
