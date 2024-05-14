package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// style
var resultSuccess = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Foreground(lipgloss.Color("#5BFF2E")).BorderForeground(lipgloss.Color("#5BFF2E")).Padding(2)
var resultFail = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Foreground(lipgloss.Color("#FF4444")).BorderForeground(lipgloss.Color("#FF4444")).Padding(2)
var commandExitHelp = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF4444")).Bold(true)
var commandHelp = lipgloss.NewStyle().Foreground(lipgloss.Color("#E449E6")).Bold(true)
var showTablesS = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Foreground(lipgloss.Color("#5BFF2E")).BorderForeground(lipgloss.Color("#5BFF2E")).Padding(1)
var skyText = lipgloss.NewStyle().Foreground(lipgloss.Color("#57FAFF")).Bold(true)

func (m model) Starter(url string) string {
	var tableString string
	app := new(App)

	if _, ok := app.PingDB(url); !ok {
		return lipgloss.Place(
			m.Width,
			m.Height,
			lipgloss.Left,
			lipgloss.Left,
			lipgloss.JoinVertical(
				lipgloss.Left,
				resultFail.Render("failed to connect to database"),
				fmt.Sprintf("Press key %s or %s to quit", commandExitHelp.Render("ctrl+c"), commandExitHelp.Render("esc")),
			))
	}

	t, ok := app.GetTables()

	if !ok {
		return lipgloss.Place(
			m.Width,
			m.Height,
			lipgloss.Left,
			lipgloss.Left,
			lipgloss.JoinVertical(
				lipgloss.Left,
				resultFail.Render("failed to get tables"),
				fmt.Sprintf("Press key %s or %s to quit", commandExitHelp.Render("ctrl+c"), commandExitHelp.Render("esc")),
			))
	}

	// create select option table
	for _, n := range t {
		m.ListOfTable = append(m.ListOfTable, TableList{name: n, checked: false})
	}

	// render m.ListOfTable
	sTable := "Tables\n"

	for i, table := range m.ListOfTable {
		cursor := " "
		// is pointed table1
		if m.Cursor == i {
			cursor = ">"
		}

		checkedStatus := "x"
		// is chacked
		if _, ok := m.Selected[i]; ok {
			table.checked = true
			tableString = table.name
			checkedStatus = "[v]"
		}

		sTable += fmt.Sprintf("%s %s %s\n", cursor, checkedStatus, table.name)
	}

	// get column tables
	columns := app.GetColumns(tableString)
	formatedColumns := strings.Join(columns, " | ")

	// get values
	values := app.GetValues(tableString)

	formatedValues := FormatValues(values, columns)

	return lipgloss.Place(
		m.Width,
		m.Height,
		lipgloss.Left,
		lipgloss.Left,
		lipgloss.JoinVertical(
			lipgloss.Left,
			resultSuccess.Render("success connect to database"),
			showTablesS.Render(sTable),
			skyText.MarginTop(1).Render(formatedColumns),
			skyText.Render(formatedValues),
			fmt.Sprintf("Press key %s or %s to quit", commandExitHelp.Render("ctrl+c"), commandExitHelp.Render("esc")),
		))
}

func FormatValues(values []map[string]interface{}, columns []string) string {
	var result strings.Builder

	for _, v := range values {
		for i, col := range columns {
			if i > 0 {
				result.WriteString(" | ")
			}
			val := v[col]
			switch t := val.(type) {
			case int:
				result.WriteString(strconv.Itoa(t))
			case string:
				result.WriteString(t)
			default:
				result.WriteString(fmt.Sprintf("%v", t))
			}
		}

		result.WriteString("\n")
	}

	return result.String()
}
