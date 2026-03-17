package main

import (
	"os"
	"time"

	tea "charm.land/bubbletea/v2"

	"studytui/modules"
	"studytui/modules/todo"
	"studytui/modules/pomodoro"
)

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(time.Time) tea.Msg {
		return pomodoro.TickMsg{}
	})
}

func main() {

	todoModule := todo.New("~/Documents/notes/tasks.json")
	pomoModule := pomodoro.New("~/Documents/notes/timerconf.json")

	manager := modules.Manager{
		Modules: []modules.Module{
			todoModule,
			pomoModule,
		},
	}

	p := tea.NewProgram(manager)

	if _, err := p.Run(); err != nil {
		println("Error:", err.Error())
		os.Exit(1)
	}
}
