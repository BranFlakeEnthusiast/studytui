package main

import (
	"os"

	tea "charm.land/bubbletea/v2"

	"studytui/modules"
	"studytui/modules/todo"
)

func main() {

	todoModule := todo.New("~/Documents/notes/tasks.json")

	manager := modules.Manager{
		Modules: []modules.Module{
			todoModule,
		},
	}

	p := tea.NewProgram(manager)

	if _, err := p.Run(); err != nil {
		println("Error:", err.Error())
		os.Exit(1)
	}
}
