package modules

import (
	"time"

	tea "charm.land/bubbletea/v2"

	"studytui/modules/pomodoro"
)

type Manager struct {
	Modules []Module
	Current int

	width int
	height int
}

func Tick() tea.Cmd {
	return tea.Tick(time.Second, func(time.Time) tea.Msg {
		return pomodoro.TickMsg{}
	})
}

func (m Manager) Init() tea.Cmd {
	return tea.Batch(
		m.Modules[m.Current].Init(),
		Tick(),
	)}

func (m Manager) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {

	case tea.KeyPressMsg:
		switch msg.String() {

		case "tab":
			m.Current++
			if m.Current >= len(m.Modules) {
				m.Current = 0
			}

			mod, _ := m.Modules[m.Current].Update(
				tea.WindowSizeMsg{
					Width: m.width,
					Height: m.height,
				},
			)

			m.Modules[m.Current] = mod.(Module)
			return m, nil

		case "shift+tab":
			m.Current--
			if m.Current < 0 {
				m.Current = len(m.Modules)-1
			}
		}

	case pomodoro.TickMsg:
		mod, cmd := m.Modules[1].Update(msg)
		m.Modules[1] = mod.(Module)
		return m, cmd

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

 	  mod, cmd := m.Modules[m.Current].Update(msg)
    m.Modules[m.Current] = mod.(Module)

		return m, cmd
	}

	module := m.Modules[m.Current]
	updated, cmd := module.Update(msg)

	m.Modules[m.Current] = updated.(Module)

	return m, cmd
}

func (m Manager) View() tea.View {
	return m.Modules[m.Current].View()
}
