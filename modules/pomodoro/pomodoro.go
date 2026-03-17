package pomodoro

import (
	"fmt"
	"time"
	"os"
	"encoding/json"
	"path/filepath"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/common-nighthawk/go-figure"
)

type model struct {
	is_break bool
	duration time.Duration
	remaining time.Duration
	running bool

	fontChoice int
	fonts []string
	path string

	width int
	height int
}

type TickMsg struct {}

type config struct {
	FontChoice int `json:"font_choice"`
}

func saveConfig(path string, c config) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func loadConfig(path string) (config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return config{}, nil
		}
		return config{}, err
	}

	var c config
	err = json.Unmarshal(data, &c)
	return c, err
}

func New(path string) model{

	if len(path) > 0 && path[:2] == "~/" {
		home, err := os.UserHomeDir()
		if err == nil {
			path = filepath.Join(home, path[2:])
		}
	}
	cfg, _ := loadConfig(path)

	duration := 25 * time.Minute

	f := []string {
"3-d",
"3x5",
"5lineoblique",
"acrobatic",
"alligator",
"alligator2",
"alphabet",
"avatar",
"banner",
"banner3-D",
"banner3",
"banner4",
"barbwire",
"basic",
"bell",
"big",
"bigchief",
"binary",
"block",
"bubble",
"bulbhead",
"calgphy2",
"caligraphy",
"catwalk",
"chunky",
"coinstak",
"colossal",
"computer",
"contessa",
"contrast",
"cosmic",
"cosmike",
"cricket",
"cursive",
"cyberlarge",
"cybermedium",
"cybersmall",
"diamond",
"digital",
"doh",
"doom",
"dotmatrix",
"drpepper",
"eftichess",
"eftifont",
"eftipiti",
"eftirobot",
"eftitalic",
"eftiwall",
"eftiwater",
"epic",
"fender",
"fourtops",
"fuzzy",
"goofy",
"gothic",
"graffiti",
"hollywood",
"invita",
"isometric1",
"isometric2",
"isometric3",
"isometric4",
"italic",
"ivrit",
"jazmine",
"jerusalem",
"katakana",
"kban",
"larry3d",
"lcd",
"lean",
"letters",
"linux",
"lockergnome",
"madrid",
"marquee",
"maxfour",
"mike",
"mini",
"mirror",
"mnemonic",
"morse",
"moscow",
"nancyj-fancy",
"nancyj-underlined",
"nancyj",
"nipples",
"ntgreek",
"o8",
"ogre",
"pawp",
"peaks",
"pebbles",
"pepper",
"poison",
"puffy",
"pyramid",
"rectangles",
"relief",
"relief2",
"rev",
"roman",
"rot13",
"rounded",
"rowancap",
"rozzo",
"runic",
"runyc",
"sblood",
"script",
"serifcap",
"shadow",
"short",
"slant",
"slide",
"slscript",
"small",
"smisome1",
"smkeyboard",
"smscript",
"smshadow",
"smslant",
"smtengwar",
"speed",
"stampatello",
"standard",
"starwars",
"stellar",
"stop",
"straight",
"tanja",
"tengwar",
"term",
"thick",
"thin",
"threepoint",
"ticks",
"ticksslant",
"tinker-toy",
"tombstone",
"trek",
"tsalagi",
"twopoint",
"univers",
"usaflag",
"wavy",
"weird",
	}

	return model{
		duration: duration,
		remaining: duration,
		fonts: f,
		fontChoice: cfg.FontChoice,
		path: path,
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

		case "f":
			m.fontChoice = (m.fontChoice + 1) % len(m.fonts)
			saveConfig(m.path, config{
				FontChoice: m.fontChoice,
			})
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

func formatTime(d time.Duration, font string) string {
	m := int(d.Minutes())
	s := int(d.Seconds()) % 60
	t := fmt.Sprintf("%02d:%02d", m, s)
	f := figure.NewFigure(t,font,true)

	return f.String()
}

var (
titleStyle = lipgloss.NewStyle().
	Bold(true).
	Align(lipgloss.Center).
	MarginBottom(1)

timerStyle = lipgloss.NewStyle().
	Bold(true)

modeStyle = lipgloss.NewStyle().
	Align(lipgloss.Center).
	MarginTop(2)

helpStyle = lipgloss.NewStyle().
	Faint(true).
	Align(lipgloss.Center)
)

func(m model) View() tea.View{
	s := titleStyle.Width(m.width).Render("Pomodoro Timer") + "\n"

	timerContainer := timerStyle.Render(formatTime(m.remaining, m.fonts[m.fontChoice]))
	s += lipgloss.NewStyle().Width(m.width).Align(lipgloss.Center).Render(timerContainer)

	if m.is_break {
		s += modeStyle.Width(m.width).Render("break!")
	} else {
		s += modeStyle.Width(m.width).Render("work!")
	}

	s += helpStyle.Width(m.width).Render("\n\n q quit - ␣ play/pause - r reset - s switch mode - f switch font")

	v:= tea.NewView(s)
	v.AltScreen = true
	return v
}

func (m model) Name() string {
	return "Pomodoro"
}
