package main

import (
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// style
var logo = lipgloss.NewStyle().Border(lipgloss.ThickBorder()).Foreground(lipgloss.Color("#F975DC")).BorderForeground(lipgloss.Color("#F975DC")).Padding(1, 2)

type QuestionStyle struct {
	BorderColor     lipgloss.Color
	InputFieldColor lipgloss.Style
}

func QuestionDefaultStyle() *QuestionStyle {
	s := new(QuestionStyle)
	s.BorderColor = lipgloss.Color("36")
	s.InputFieldColor = lipgloss.NewStyle().
		BorderForeground(s.BorderColor).
		BorderStyle(lipgloss.NormalBorder()).
		Padding(2).
		Width(80)

	return s
}

// we working in here
type model struct {
	Index              int
	Width              int
	Height             int
	QuestionCredential []QuestionCredential
	Done               bool
	ListOfTable        []TableList
	Cursor             int
	Selected           map[int]struct{}

	// style
	QuestionStyle *QuestionStyle
}

type QuestionCredential struct {
	index    string
	question string
	answer   string
	input    Input
}

type TableList struct {
	name    string
	checked bool
}

func NewQuestion(index string, question string) QuestionCredential {
	return QuestionCredential{index: index, question: question}
}

func newShortQuestion(index string, question string) QuestionCredential {
	q := NewQuestion(index, question)
	field := NewShortQuestionField()
	q.input = field
	return q
}

func QuestionModel(credential []QuestionCredential) *model {
	// question style
	styleQ := QuestionDefaultStyle()

	clientField := textinput.New()
	clientField.Placeholder = "Enter here"
	clientField.Focus()

	return &model{
		QuestionCredential: credential,
		QuestionStyle:      styleQ,
		Selected:           make(map[int]struct{}),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	current := &m.QuestionCredential[m.Index]

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
	case tea.KeyMsg:
		if !m.Done {
			switch msg.String() {
			case "ctrl+c", "esc":
				return m, tea.Quit
			case "enter":
				if m.Index == len(m.QuestionCredential)-1 {
					m.Done = true
				}
				current.answer = current.input.Value()
				m.Next()
				return m, current.input.Blur
			}
		} else {
			switch msg.String() {
			case "ctrl+c", "esc":
				return m, tea.Quit
			case "up":
				if m.Cursor > 0 {
					m.Cursor--
				}
			case "down":
				m.Cursor++
			case "enter":
				_, ok := m.Selected[m.Cursor]
				if ok {
					delete(m.Selected, m.Cursor)
				} else {
					m.Selected[m.Cursor] = struct{}{}
				}
			}
		}
	}

	current.input, cmd = current.input.Update(msg)

	return m, cmd
}

func (m model) View() string {
	current := &m.QuestionCredential[m.Index]

	// if done print here
	if m.Done {
		var username string
		var password string
		var dbname string
		var url string
		for _, q := range m.QuestionCredential {
			if q.index == "username" {
				username = q.answer
			} else if q.index == "password" {
				password = q.answer
			} else if q.index == "dbname" {
				dbname = q.answer
			}
		}

		url = fmt.Sprintf("%s%s@tcp(127.0.0.1:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, dbname)

		// test connection
		output := m.Starter(url)

		return output

	}

	if m.Width == 0 {
		return fmt.Sprintf("%s\n\nLoading...\n%s", logo.Render("DBPull"), "")
	}

	return lipgloss.Place(
		m.Width,
		m.Height,
		lipgloss.Left,
		lipgloss.Left,
		lipgloss.JoinVertical(lipgloss.Left,
			logo.Render("DBPull\n"),
			m.QuestionCredential[m.Index].question,
			m.QuestionStyle.InputFieldColor.Render(current.input.View()),
			fmt.Sprintf("Press key %s or %s to quit", commandHelp.Render("ctrl+c"), commandHelp.Render("esc")),
		),
	)

}

func (m *model) Next() {
	if m.Index < len(m.QuestionCredential)-1 {
		m.Index++
	} else {
		m.Index = 0
	}
}

func main() {

	questionCredential := []QuestionCredential{newShortQuestion("username", "username?"), newShortQuestion("password", "password example: ':mypassword' if password null its ok leave with null value"), newShortQuestion("dbname", "DB name?")}

	m := QuestionModel(questionCredential)

	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	app := tea.NewProgram(m)

	if _, err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
