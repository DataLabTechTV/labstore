package render

import (
	"os"
	"strings"

	"github.com/mattn/go-isatty"
)

func Box(lines ...string) string {
	var box strings.Builder

	width := 0

	for _, line := range lines {
		if length := stringDisplayWidth(line); length > width {
			width = length
		}
	}

	border := strings.Repeat("─", width+2)
	topBorder := "╭" + border + "╮"
	bottomBorder := "╰" + border + "╯"

	box.WriteString(topBorder)
	box.WriteRune('\n')

	for _, line := range lines {
		row := "│ " + line

		if pad := width - stringDisplayWidth(line); pad > 0 {
			row += strings.Repeat(" ", pad)
		}

		row += " │"
		box.WriteString(row)
		box.WriteRune('\n')
	}

	box.WriteString(bottomBorder)

	return box.String()
}

func SupportsBox() bool {
	return isatty.IsTerminal(os.Stdout.Fd())
}

func runeDisplayWidth(r rune) int {
	if r <= 0x1FFF {
		return 1
	}
	return 2
}

func stringDisplayWidth(s string) int {
	width := 0
	for _, r := range s {
		width += runeDisplayWidth(r)
	}
	return width
}
