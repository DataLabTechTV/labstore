package infopane

const ValueNone = "<none>"

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
