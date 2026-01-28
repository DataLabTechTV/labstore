package render

import "github.com/charmbracelet/lipgloss"

type Palette struct {
	TextPrimary  lipgloss.Color
	TextMuted    lipgloss.Color
	TextInverted lipgloss.Color

	Surface      lipgloss.Color
	SurfaceAlt   lipgloss.Color
	SurfaceHover lipgloss.Color

	Accent      lipgloss.Color
	AccentMuted lipgloss.Color

	Success lipgloss.Color
	Warning lipgloss.Color
	Error   lipgloss.Color

	Border lipgloss.Color
}

var (
	ActivePalette = DefaultPalette

	DefaultPalette = Palette{
		TextPrimary:  lipgloss.Color("014"),
		TextMuted:    lipgloss.Color("008"),
		TextInverted: lipgloss.Color("005"),

		Surface:      lipgloss.Color("000"),
		SurfaceAlt:   lipgloss.Color("008"),
		SurfaceHover: lipgloss.Color("007"),

		Accent:      lipgloss.Color("003"),
		AccentMuted: lipgloss.Color("012"),

		Success: lipgloss.Color("010"),
		Warning: lipgloss.Color("011"),
		Error:   lipgloss.Color("009"),

		Border: lipgloss.Color("015"),
	}
)

func Swatch(label string, fg, bg *lipgloss.Color) string {
	style := lipgloss.NewStyle()

	if fg != nil {
		style = style.Foreground(fg)
	}

	if bg != nil {
		style = style.Background(bg)
	}

	return style.Render(label)
}

func (p *Palette) Render() string {
	h1 := lipgloss.NewStyle().
		Bold(true).
		Foreground(p.Accent).
		Underline(true)

	left := lipgloss.NewStyle().Width(25).Align(lipgloss.Left)
	right := lipgloss.NewStyle().Width(30).Align(lipgloss.Right)

	rows := []string{
		h1.Render("Color Palette\n"),

		lipgloss.JoinHorizontal(
			lipgloss.Top,
			left.Render(Swatch("TextPrimary", &p.TextPrimary, nil)),
			right.Render(Swatch("TextPrimary on Surface", &p.TextPrimary, &p.Surface)),
		),
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			left.Render(Swatch("TextMuted", &p.TextMuted, nil)),
			right.Render(Swatch("TextMuted on Surface", &p.TextMuted, &p.Surface)),
		),
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			left.Render(Swatch("TextInverted", &p.TextInverted, nil)),
			right.Render(Swatch("TextInverted on Accent", &p.TextInverted, &p.Accent)),
		),

		lipgloss.JoinHorizontal(
			lipgloss.Top,
			left.Render(Swatch("Error", &p.Error, nil)),
			right.Render(Swatch("TextInverted on Error", &p.TextInverted, &p.Error)),
		),
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			left.Render(Swatch("Warning", &p.Warning, nil)),
			right.Render(Swatch("TextInverted on Warning", &p.TextInverted, &p.Warning)),
		),
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			left.Render(Swatch("Success", &p.Success, nil)),
			right.Render(Swatch("TextInverted on Success", &p.TextInverted, &p.Success)),
		),

		lipgloss.JoinHorizontal(
			lipgloss.Top,
			left.Render(Swatch("Accent on Surface", &p.Accent, &p.Surface)),
			right.Render(Swatch("TextPrimary on SurfaceAlt", &p.TextPrimary, &p.SurfaceAlt)),
		),
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			left.Render(Swatch("AccentMuted on Surface", &p.AccentMuted, &p.Surface)),
			right.Render(Swatch("TextInverted on SurfaceHover", &p.TextInverted, &p.SurfaceHover)),
		),
	}

	boxedStyle := lipgloss.NewStyle().
		Padding(0, 1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(p.Border)

	return boxedStyle.Render(lipgloss.JoinVertical(lipgloss.Left, rows...))
}
