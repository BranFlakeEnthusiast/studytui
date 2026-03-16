package modules

import tea "charm.land/bubbletea/v2"

type Manager struct {
	Modules []Module
	Current int
}

func (m Manager) Init() tea.Cmd {
	return m.Modules[m.Current].Init()
}

func (m Manager) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {

	case tea.KeyPressMsg:
		switch msg.String() {

		case "tab":
			m.Current++
			if m.Current >= len(m.Modules) {
				m.Current = 0
			}
			return m, nil

		case "shift + tab":
			m.Current--
			if m.Current < 0 {
				m.Current = len(m.Modules)
			}
			return m, nil
		}
	}

	module := m.Modules[m.Current]
	updated, cmd := module.Update(msg)

	m.Modules[m.Current] = updated.(Module)

	return m, cmd
}

func (m Manager) View() tea.View {
	return m.Modules[m.Current].View()
}
