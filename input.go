package main

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Input interface {
	Value() string
	Blur() tea.Msg
	Update(tea.Msg) (Input, tea.Cmd)
	View() string
}

// text field
type ShortAnswerField struct {
	textfield textinput.Model
}

func NewShortQuestionField() *ShortAnswerField {
	tf := textinput.New()
	tf.Placeholder = "Enter here"
	tf.Focus()

	return &ShortAnswerField{textfield: tf}
}

func (sa *ShortAnswerField) Value() string {
	return sa.textfield.Value()
}

func (sa *ShortAnswerField) Blur() tea.Msg {
	return sa.textfield.Blur
}

func (sa *ShortAnswerField) Update(msg tea.Msg) (Input, tea.Cmd) {
	var cmd tea.Cmd
	sa.textfield, cmd = sa.textfield.Update(msg)

	return sa, cmd
}

func (sa *ShortAnswerField) View() string {
	return sa.textfield.View()
}

// textarea
// type LongAnswerField struct {
// 	textarea textarea.Model
// }

// func NewLongQuestionField() *LongAnswerField {
// 	ta := textarea.New()
// 	ta.Placeholder = "Enter here"
// 	ta.Focus()

// 	return &LongAnswerField{textarea: ta}
// }

// func (la *LongAnswerField) Value() string {
// 	return la.textarea.Value()
// }

// func (la *LongAnswerField) Blur() tea.Msg {
// 	return la.textarea.Blur
// }

// func (la *LongAnswerField) Update(msg tea.Msg) (Input, tea.Cmd) {
// 	var cmd tea.Cmd
// 	la.textarea, cmd = la.textarea.Update(msg)

// 	return la, cmd
// }

// func (la *LongAnswerField) View() string {
// 	return la.textarea.View()
// }
