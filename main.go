package main

import (
	"os"
	"time"

	tea "charm.land/bubbletea/v2"

	"studytui/modules"
	"studytui/modules/flashcard"
	"studytui/modules/pomodoro"
	"studytui/modules/todo"
)

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(time.Time) tea.Msg {
		return pomodoro.TickMsg{}
	})
}

func main() {

	todoModule := todo.New("~/Documents/notes/tasks.json")
	pomoModule := pomodoro.New("~/Documents/notes/timerconf.json")
	flashModule := flashcard.New("~/Documents/notes/flashcards.json")

	manager := modules.Manager{
		Modules: []modules.Module{
			todoModule,
			pomoModule,
			flashModule,
		},
	}

	p := tea.NewProgram(manager)

	if _, err := p.Run(); err != nil {
		println("Error:", err.Error())
		os.Exit(1)
	}
}
