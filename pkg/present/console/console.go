package console

import (
	"strconv"
	"strings"

	"github.com/rs/zerolog"

	"github.com/agrski/greg/pkg/match"
	"github.com/agrski/greg/pkg/types"
)

type Console struct {
	enableColour bool
	logger       zerolog.Logger
}

func New(logger zerolog.Logger, enableColour bool) *Console {
	return &Console{
		enableColour: enableColour,
		logger:       logger,
	}
}

func (c *Console) Write(fileInfo *types.FileInfo, match *match.Match) {
	c.logger.Log().Msg(string(fgBlue) + fileInfo.Path + string(reset))

	for _, p := range match.Positions {
		sb := strings.Builder{}

		line := strconv.Itoa(int(p.Line + 1))

		if c.enableColour {
			// Line number
			sb.WriteString(string(fgMagenta))
			sb.WriteString(line)
			sb.WriteString(string(reset))
			sb.WriteByte(':')
			// Text
			sb.WriteString(p.Text[:p.ColumnStart])
			sb.WriteString(string(fgRed))
			sb.WriteString(p.Text[p.ColumnStart:p.ColumnEnd])
			sb.WriteString(string(reset))
			sb.WriteString(p.Text[p.ColumnEnd:])
		} else {
			// Line number
			sb.WriteString(line)
			sb.WriteByte(':')
			// Text
			sb.WriteString(p.Text)
		}

		c.logger.Log().Msg(sb.String())
	}

	c.logger.Log().Send()
}
