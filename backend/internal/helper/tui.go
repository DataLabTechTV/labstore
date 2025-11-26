package helper

import (
	"fmt"
	"strings"
)

func RuneDisplayWidth(r rune) int {
	if r <= 0x1FFF {
		return 1
	}
	return 2
}

func StringDisplayWidth(s string) int {
	width := 0
	for _, r := range s {
		width += RuneDisplayWidth(r)
	}
	return width
}

func Box(lines ...string) {
	width := 0

	for _, line := range lines {
		if length := StringDisplayWidth(line); length > width {
			width = length
		}
	}

	border := strings.Repeat("─", width+2)
	topBorder := "╭" + border + "╮"
	bottomBorder := "╰" + border + "╯"

	fmt.Println(topBorder)

	for _, line := range lines {
		row := "│ " + line

		if pad := width - StringDisplayWidth(line); pad > 0 {
			row += strings.Repeat(" ", pad)
		}

		row += " │"
		fmt.Println(row)
	}

	fmt.Println(bottomBorder)
}
