package console

import (
	"io"
	"strconv"
	"strings"

	"github.com/agrski/greg/pkg/match"
	"github.com/agrski/greg/pkg/types"
)

type Console struct {
	enableColour bool
	out          io.StringWriter
}

func New(out io.StringWriter, enableColour bool) *Console {
	return &Console{
		enableColour: enableColour,
		out:          out,
	}
}

func (c *Console) Write(fileInfo *types.FileInfo, match *match.Match) {
	_, err := c.out.WriteString(string(fgBlue) + fileInfo.Path + string(reset) + "\n")
	if err != nil {
		return
	}

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

		sb.WriteString("\n")

		_, err := c.out.WriteString(sb.String())
		if err != nil {
			return
		}
	}

	_, _ = c.out.WriteString("\n")
}
