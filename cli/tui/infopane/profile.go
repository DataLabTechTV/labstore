package infopane

type ProfileInfoPane struct {
	Model
}

func NewProfile(label, value string) ProfileInfoPane {
	return ProfileInfoPane{Model: New(label, value)}
}

func (m ProfileInfoPane) Clear() ProfileInfoPane {
	m.Model = m.Model.Clear()
	return m
}
