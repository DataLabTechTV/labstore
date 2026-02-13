package infopane

const ValueNone = ""

type Model struct {
	Label  string
	Value  string
	Width  int
	Height int
}

func New(label string, value string) Model {
	return Model{
		Label:  label,
		Value:  value,
		Height: 3,
	}
}

func (m *Model) Clear() {
	m.Value = ""
}
