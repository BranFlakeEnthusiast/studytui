package pomodoro

import (
	"fmt"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/common-nighthawk/go-figure"
)

type model struct {
	is_break bool
	duration time.Duration
	remaining time.Duration
	running bool

	width int
	height int
}

type TickMsg struct {}

func New() model{
	duration := 25 * time.Minute

	return model{
		duration: duration,
		remaining: duration,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func Tick() tea.Cmd {
	return tea.Tick(time.Second, func(time.Time) tea.Msg {
		return TickMsg{}
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyPressMsg:

		switch msg.String(){

		case "q","ctrl+c":
			return m, tea.Quit

		case "space":
			m.running = !m.running

		case "r":
			m.running = false
			m.remaining = m.duration

		case "s":
			m.running = false
			m.is_break = !m.is_break

			if m.is_break{
				m.duration = 5 * time.Minute
			} else {
				m.duration = 25 * time.Minute
			}

			m.remaining = m.duration
		}

		case TickMsg:
			if m.running{
				m.remaining -= time.Second

				if m.remaining <= 0 {
					m.running = false
					m.is_break = !m.is_break

				if m.is_break{
					m.duration = 5 * time.Minute
				} else {
					m.duration = 25 * time.Minute
				}

				m.remaining = m.duration
				}
			}
			return m, Tick()
	}

return m, nil
}

func formatTime(d time.Duration) string {
	m := int(d.Minutes())
	s := int(d.Seconds()) % 60
	t := fmt.Sprintf("%02d:%02d", m, s)
	f := figure.NewFigure(t,"hollywood",true)

	return f.String()
}

var (
titleStyle = lipgloss.NewStyle().
	Bold(true).
	Align(lipgloss.Center).
	MarginBottom(1)

timerStyle = lipgloss.NewStyle().
	Bold(true)

helpStyle = lipgloss.NewStyle().
	Faint(true).
	Align(lipgloss.Center)
)

func(m model) View() tea.View{
	s := titleStyle.Width(m.width).Render("Pomodoro Timer") + "\n"

	timerContainer := timerStyle.Render(formatTime(m.remaining))
	s += lipgloss.NewStyle().Width(m.width).Align(lipgloss.Center).Render(timerContainer)

	s += helpStyle.Width(m.width).Render("\n\n q quit - ␣ play/pause - r reset - s switch mode")

	v:= tea.NewView(s)
	v.AltScreen = true
	return v
}

func (m model) Name() string {
	return "Pomodoro"
}
