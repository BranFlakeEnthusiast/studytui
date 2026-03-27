package flashcard

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type Card struct {
	Face string
	Back string
}

type addStep int

const (
	addStepFace addStep = iota
	addStepBack
)

type model struct {
	Cards   []Card
	Flipped bool
	Cursor  int

	Input    textinput.Model
	Adding   bool
	addStep  addStep
	pendingFace string

	Editing     bool
	editingBack bool

	Path   string
	width  int
	height int
}

func New(path string) model {
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

	cards, err := loadCards(path)
	if err != nil {
		cards = []Card{}
	}

	return model{
		Path:  path,
		Cards: cards,
		Input: ti,
	}
}

func saveCards(path string, cards []Card) error {
	data, err := json.MarshalIndent(cards, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func loadCards(path string) ([]Card, error) {
	os.MkdirAll(filepath.Dir(path), 0755)

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return []Card{}, nil
		}
		return nil, err
	}

	var cards []Card
	err = json.Unmarshal(data, &cards)
	if err != nil {
		return nil, err
	}
	return cards, nil
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyPressMsg:

		if m.Adding {
			var cmd tea.Cmd
			m.Input, cmd = m.Input.Update(msg)

			switch msg.String() {
			case "enter":
				val := m.Input.Value()
				if m.addStep == addStepFace {
					if val != "" {
						m.pendingFace = val
						m.Input.SetValue("")
						m.Input.Placeholder = "Back of card..."
						m.addStep = addStepBack
					}
				} else {
					// back step — val may be empty (blank back is fine)
					m.Cards = append(m.Cards, Card{Face: m.pendingFace, Back: val})
					m.Cursor = len(m.Cards) - 1
					saveCards(m.Path, m.Cards)
					m.Input.SetValue("")
					m.Input.Placeholder = ""
					m.pendingFace = ""
					m.addStep = addStepFace
					m.Adding = false
					m.Flipped = false
				}

			case "esc":
				m.Input.SetValue("")
				m.Input.Placeholder = ""
				m.pendingFace = ""
				m.addStep = addStepFace
				m.Adding = false
			}
			return m, cmd
		}

		if m.Editing {
			var cmd tea.Cmd
			m.Input, cmd = m.Input.Update(msg)

			switch msg.String() {
			case "enter":
				val := m.Input.Value()
				if !m.editingBack {
					// finished editing face
					if val != "" {
						m.Cards[m.Cursor].Face = val
					}
					m.Input.SetValue(m.Cards[m.Cursor].Back)
					m.Input.Placeholder = "Back of card..."
					m.editingBack = true
				} else {
					// finished editing back
					m.Cards[m.Cursor].Back = val
					saveCards(m.Path, m.Cards)
					m.Input.SetValue("")
					m.Input.Placeholder = ""
					m.editingBack = false
					m.Editing = false
					m.Flipped = false
				}

			case "esc":
				m.Input.SetValue("")
				m.Input.Placeholder = ""
				m.editingBack = false
				m.Editing = false
			}
			return m, cmd
		}

		switch msg.String() {
		case "q", "ctrl+c":
			saveCards(m.Path, m.Cards)
			return m, tea.Quit

		case "l", "right":
			if len(m.Cards) > 0 && m.Cursor < len(m.Cards)-1 {
				m.Cursor++
				m.Flipped = false
			}

		case "h", "left":
			if m.Cursor > 0 {
				m.Cursor--
				m.Flipped = false
			}

		case "space", "enter", "f":
			if len(m.Cards) > 0 {
				m.Flipped = !m.Flipped
			}

		case "a":
			m.Input.Placeholder = "Front of card..."
			m.addStep = addStepFace
			m.Adding = true

		case "e":
			if len(m.Cards) > 0 {
				m.Input.SetValue(m.Cards[m.Cursor].Face)
				m.Input.Placeholder = "Front of card..."
				m.editingBack = false
				m.Editing = true
				m.Flipped = false
			}

		case "d", "backspace":
			if len(m.Cards) > 0 {
				m.Cards = append(m.Cards[:m.Cursor], m.Cards[m.Cursor+1:]...)
				if m.Cursor > 0 {
					m.Cursor--
				}
				saveCards(m.Path, m.Cards)
				m.Flipped = false
			}
		}
	}

	return m, nil
}


var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Align(lipgloss.Center).
			MarginBottom(1)

	cardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(1, 4).
			Align(lipgloss.Center)

	faceStyle = lipgloss.NewStyle().
			Bold(true).
			MarginBottom(2)

	backStyle = lipgloss.NewStyle().
			MarginBottom(2)
	labelStyle = lipgloss.NewStyle().
			Faint(true).
			Italic(true)

	counterStyle = lipgloss.NewStyle().
			Faint(true).
			Align(lipgloss.Center)

	emptyStyle = lipgloss.NewStyle().
			Faint(true).
			Align(lipgloss.Center)

	helpStyle = lipgloss.NewStyle().
			Faint(true).
			Align(lipgloss.Center)
)

func (m model) View() tea.View {
	s := titleStyle.Width(m.width).Render("Flashcards") + "\n"

	if len(m.Cards) == 0 {
		s += emptyStyle.Width(m.width).Render("\n\nNo cards yet. Press a to add one.\n\n")
	} else {
		card := m.Cards[m.Cursor]

		var cardContent string
		if m.Flipped {
			cardContent = labelStyle.Render("back") + "\n\n\n" + backStyle.Render(card.Back)
		} else {
			cardContent = labelStyle.Render("front") + "\n\n\n" + faceStyle.Render(card.Face)
		}

		cardWidth := min(m.width-4, 60)
		rendered := cardStyle.Width(cardWidth).Render(cardContent)
		s += lipgloss.NewStyle().Width(m.width).Align(lipgloss.Center).Render(rendered) + "\n"

		counter := fmt.Sprintf("%d / %d", m.Cursor+1, len(m.Cards))
		s += counterStyle.Width(m.width).Render(counter)
	}

	if m.Adding || m.Editing {
		prompt := "Front: "
		if m.addStep == addStepBack || m.editingBack {
			prompt = "Back:  "
		}
		s += "\n\n" + prompt + m.Input.View()
	} else {
		s += helpStyle.Width(m.width).Render("\n\n q quit · a add · e edit · d delete · ←/→ navigate · ␣ flip")
	}

	v := tea.NewView(s)
	v.AltScreen = true
	return v
}

func (m model) Name() string {
	return "Flashcards"
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
