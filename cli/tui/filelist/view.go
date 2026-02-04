package filelist

func (m Model) View() string {
	if len(m.Entries) < 1 {
		return ""
	}

	return m.table.View()
}
