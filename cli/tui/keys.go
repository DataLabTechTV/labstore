package tui

import "github.com/charmbracelet/bubbles/key"

type KeyMap interface {
	HelpKeys() []key.Binding
}

type HomeKeyMap struct {
	Quit     key.Binding
	Profiles key.Binding
	Put      key.Binding
	Get      key.Binding
	Delete   key.Binding
	Head     key.Binding
	NavUp    key.Binding
	Open     key.Binding
	Select   key.Binding
	Next     key.Binding
	Previous key.Binding

	Down key.Binding
	Up   key.Binding
	End  key.Binding
	Home key.Binding
	PgDn key.Binding
	PgUp key.Binding

	Focus1 key.Binding
	Focus2 key.Binding
	Focus3 key.Binding
	Focus4 key.Binding
}

func (k HomeKeyMap) HelpKeys() []key.Binding {
	return []key.Binding{
		k.Quit,
		k.Profiles,
		k.Put,
		k.Get,
		k.Delete,
		k.Head,
		k.NavUp,
		k.Open,
		k.Select,
		k.Next,
		k.Previous,
	}
}

var DefaultHomeKeyMap KeyMap = HomeKeyMap{
	Profiles: key.NewBinding(key.WithKeys("ctrl+p"), key.WithHelp("ctrl+p", "Profiles")),
	Put:      key.NewBinding(key.WithKeys("p"), key.WithHelp("p", "PUT")),
	Get:      key.NewBinding(key.WithKeys("g"), key.WithHelp("g", "GET")),
	Delete:   key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "DELETE")),
	Head:     key.NewBinding(key.WithKeys("h"), key.WithHelp("h", "HEAD")),
	Open:     key.NewBinding(key.WithKeys("right", "enter"), key.WithHelp("→/⏎", "Open")),
	NavUp:    key.NewBinding(key.WithKeys("left", "backspace"), key.WithHelp("←/⌫ ", "Up a Level")),
	Select:   key.NewBinding(key.WithKeys("space"), key.WithHelp("␣", "Select")),
	Next:     key.NewBinding(key.WithKeys("tab"), key.WithHelp("<tab>", "Next")),
	Previous: key.NewBinding(key.WithKeys("shift+tab"), key.WithHelp("shift+<tab>", "Previous")),
	Quit:     key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "Quit")),

	Down: key.NewBinding(key.WithKeys("down", "j")),
	Up:   key.NewBinding(key.WithKeys("up", "k")),
	End:  key.NewBinding(key.WithKeys("end", "$")),
	Home: key.NewBinding(key.WithKeys("home", "0")),
	PgDn: key.NewBinding(key.WithKeys("pgdown")),
	PgUp: key.NewBinding(key.WithKeys("pgup")),

	Focus1: key.NewBinding(key.WithKeys("1")),
	Focus2: key.NewBinding(key.WithKeys("2")),
	Focus3: key.NewBinding(key.WithKeys("3")),
	Focus4: key.NewBinding(key.WithKeys("4")),
}
