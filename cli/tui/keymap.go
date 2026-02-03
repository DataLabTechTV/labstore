package tui

import "github.com/charmbracelet/bubbles/key"

type KeyMap interface {
	All() []key.Binding
}

type HomeKeyMap struct {
	Profiles key.Binding
	Put      key.Binding
	Get      key.Binding
	Delete   key.Binding
	Head     key.Binding
	Select   key.Binding
	Next     key.Binding
	Previous key.Binding
	Quit     key.Binding
	Focus1   key.Binding
	Focus2   key.Binding
	Focus3   key.Binding
	Focus4   key.Binding
}

func (k HomeKeyMap) All() []key.Binding {
	return []key.Binding{k.Profiles, k.Put, k.Get, k.Delete, k.Head, k.Select, k.Next, k.Previous}
}

var DefaultKeyMap KeyMap = HomeKeyMap{
	Profiles: key.NewBinding(key.WithKeys("ctrl+p"), key.WithHelp("ctrl+p", "Profiles")),
	Put:      key.NewBinding(key.WithKeys("p"), key.WithHelp("p", "PUT")),
	Get:      key.NewBinding(key.WithKeys("g"), key.WithHelp("g", "GET")),
	Delete:   key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "DELETE")),
	Head:     key.NewBinding(key.WithKeys("h"), key.WithHelp("h", "HEAD")),
	Select:   key.NewBinding(key.WithKeys("space"), key.WithHelp("<space>", "Select")),
	Next:     key.NewBinding(key.WithKeys("tab"), key.WithHelp("<tab>", "Next")),
	Previous: key.NewBinding(key.WithKeys("shift+tab"), key.WithHelp("shift+<tab>", "Previous")),
	Quit:     key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "Quit")),
	Focus1:   key.NewBinding(key.WithKeys("1")),
	Focus2:   key.NewBinding(key.WithKeys("2")),
	Focus3:   key.NewBinding(key.WithKeys("3")),
	Focus4:   key.NewBinding(key.WithKeys("4")),
}
