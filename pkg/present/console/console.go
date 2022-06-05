package console

import (
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
	c.logger.Log().Msg(string(FgBlue) + fileInfo.Path + string(Reset))

	for _, p := range match.Positions {
		sb := strings.Builder{}

		if c.enableColour {
			// Line number
			sb.WriteString(string(FgMagenta))
			sb.WriteString(fileInfo.Path)
			sb.WriteString(string(Reset))
			sb.WriteByte(':')
			// Text
			sb.WriteString(p.Text[:p.ColumnStart])
			sb.WriteString(string(FgRed))
			sb.WriteString(p.Text[p.ColumnStart:p.ColumnEnd])
			sb.WriteString(string(Reset))
			sb.WriteString(p.Text[p.ColumnEnd:])
		} else {
			// Line number
			sb.WriteString(fileInfo.Path)
			sb.WriteByte(':')
			// Text
			sb.WriteString(p.Text)
		}

		c.logger.Log().Msg(sb.String())
	}

	c.logger.Log().Send()
}
