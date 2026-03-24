package todo

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type Task struct{
	Title string
	Completed bool
}

type model struct {
	Tasks []Task
	Cursor int
	offset int

	Input textinput.Model
	Adding bool
	Editing bool

	Path string

	width int
	height int
}

func New(path string) model{
	ti := textinput.New()
	ti.Focus()
	ti.CharLimit = 200
	ti.SetWidth(50)

	if len(path) > 0 && path[:2] == "~/" {
		home, err := os.UserHomeDir()
		if err == nil {
			path = filepath.Join(home, path[2:])
		}
	}
	tasks, err := loadTasks(path)
	if err != nil{
		tasks = []Task{}
	}
	return model{
		Path: path,
		Tasks: tasks,
		Input: ti,
	}
}

func saveTasks(path string, tasks []Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func loadTasks(path string) ([]Task, error) {
	os.MkdirAll(filepath.Dir(path), 0755)

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return []Task{}, nil
		}
		return nil, err
	}

	var tasks []Task
	err = json.Unmarshal(data, &tasks)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type){

	case tea.KeyPressMsg:

		if m.Adding {
			var cmd tea.Cmd
			m.Input, cmd = m.Input.Update(msg)

			switch msg.String(){

			case "enter":
				title := m.Input.Value()

				if title != ""{
					m.Tasks = append(m.Tasks, Task{Title: title})
					saveTasks(m.Path, m.Tasks)
				}

				m.Input.SetValue("")
				m.Adding = false

			case "esc":
				m.Input.SetValue("")
				m.Adding = false
			}
			return m, cmd
  	}

		if m.Editing {
			var cmd tea.Cmd
			m.Input, cmd = m.Input.Update(msg)

			switch msg.String(){

			case "enter":
				title := m.Input.Value()

				if title != ""{
					m.Tasks[m.Cursor].Title = title
					saveTasks(m.Path, m.Tasks)
				}

				m.Input.SetValue("")
				m.Editing = false

			case "esc":
				m.Input.SetValue("")
				m.Editing = false
			}
			return m, cmd
		}

		switch msg.String(){
		case "q":
			saveTasks(m.Path, m.Tasks)
			return m, tea.Quit
		case "j", "down":
			if m.Cursor < len(m.Tasks)-1{
				m.Cursor++
			}

			visibleHeight := m.height - 6
			if m.Cursor >= m.offset + visibleHeight {
				m.offset++
			}

		case "k", "up":
			if m.Cursor > 0 {
				m.Cursor--
			}

			if m.Cursor < m.offset {
				m.offset--
			}

		case "x","enter", "space":
			if len(m.Tasks) > 0 {
				m.Tasks[m.Cursor].Completed = !m.Tasks[m.Cursor].Completed
				saveTasks(m.Path, m.Tasks)
			}

		case "d", "backspace":
			if len(m.Tasks) > 0 {
				m.Tasks = append(m.Tasks[:m.Cursor], m.Tasks[m.Cursor+1:]...,)
				if m.Cursor > 0 {
					m.Cursor--
				}
				visibleHeight := m.height - 6
				if m.offset > len(m.Tasks) - visibleHeight {
					m.offset = max(0, len(m.Tasks) - visibleHeight)
				}
				saveTasks(m.Path, m.Tasks)
			}

		case "D":
			if len(m.Tasks) > 0 {
				for i := len(m.Tasks) - 1; i >= 0; i-- {
					if m.Tasks[i].Completed{
						m.Tasks = append(m.Tasks[:i], m.Tasks[i+1:]...,)
					}
				}
				saveTasks(m.Path, m.Tasks)

				if m.Cursor > len(m.Tasks)-1{
					m.Cursor = len(m.Tasks)-1
				}
			}

		case "a":
			m.Input.Placeholder = "New task... "
			m.Adding = true

		case "e":
			if len(m.Tasks) > 0{
				m.Input.SetValue(m.Tasks[m.Cursor].Title)
				m.Editing = true
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

return m, nil
}

var (
titleStyle = lipgloss.NewStyle().
	Bold(true).
	Align(lipgloss.Center).
	MarginBottom(1)

doneStyle = lipgloss.NewStyle().
	Strikethrough(true).
	Faint(true)

checkStyle = lipgloss.NewStyle().
	Faint(true)

helpStyle = lipgloss.NewStyle().
	Faint(true).
	Align(lipgloss.Center)

cursorStyle = lipgloss.NewStyle().
	Bold(true)

taskStyle = lipgloss.NewStyle().
	Align(lipgloss.Left)
)


func (m model) View() tea.View {
	s:= titleStyle.Width(m.width).Render("Todo List") + "\n"

	visibleHeight := m.height - 6

	tasks := ""
	for i := m.offset; i < len(m.Tasks) && i < m.offset + visibleHeight; i++{

		task := m.Tasks[i]

		cursor := " "
		if m.Cursor == i {
			cursor = cursorStyle.Render(">")
		}

		check := "[ ]"
		if task.Completed {
			check = checkStyle.Render("[X]")
		}

		title := task.Title
		if task.Completed{
			title = doneStyle.Render(title)
		}

		tasks += fmt.Sprintf("%s %s %s\n", cursor, check, title)
	}
	taskContaier := taskStyle.Render(tasks)
	s += lipgloss.NewStyle().Width(m.width).Align(lipgloss.Center).Render(taskContaier)

	if m.Adding ||  m.Editing {
		s +="\n\n"+ m.Input.View()
	} else {
		s += helpStyle.Width(m.width).Render("\n\n q quit - a add - ␣ toggle - d delete - D delete checked - e edit")
	}

	v := tea.NewView(s)
	v.AltScreen = true
	return v
}

func (m model) Name() string {
	return "Todo"
}
