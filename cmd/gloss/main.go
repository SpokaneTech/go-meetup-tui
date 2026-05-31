package main

import (
	"log"
	"strings"

	tea "github.com/charmbracelet/bubbletea" // obnoxious!
	"github.com/charmbracelet/lipgloss"
)

func main() {
	m := &Model{}

	_, err := tea.NewProgram(m).Run()
	if err != nil {
		log.Fatalf("error with program: %s", err.Error())
	}
}

// Model is the main container for our tui.
// It must satisfy the tea `model` interface
type Model struct {
	bold, italic, strike, underline, background bool
}

// Init the tui. The returned tea.Cmd represents any initial
// I/O or setup type functions we want to execute
func (m Model) Init() tea.Cmd {
	return nil
}

// Update the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()
		switch key {
		case "ctrl+c", "q":
			// return the model, and ask for tea to quit
			return m, tea.Quit
		case "b":
			m.bold = !m.bold
			return m, nil
		case "s":
			m.strike = !m.strike
			return m, nil
		case "i":
			m.italic = !m.italic
			return m, nil
		case "u":
			m.underline = !m.underline
			return m, nil
		case "B":
			m.background = !m.background
			return m, nil
		}
	}

	return m, nil
}

// Our styles

var (
	highlightColor = lipgloss.AdaptiveColor{Light: "#5bfa11", Dark: "#5fdb25"}
	textColor      = lipgloss.AdaptiveColor{Light: "#0d0000", Dark: "#fcfcfc"}
	bgColor        = lipgloss.Color("#87b074")
	fgColor        = lipgloss.Color("#0d0000")
	windowStyle    = lipgloss.NewStyle().
			BorderForeground(highlightColor).
			Foreground(textColor).
			Padding(2, 0).
			Align(lipgloss.Center).
			Border(lipgloss.NormalBorder()).
			Width(40).
			Height(20)
	helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
)

func (m Model) View() string {
	doc := strings.Builder{}

	renderStyle := windowStyle.Bold(m.bold).Strikethrough(m.strike).Italic(m.italic).Underline(m.underline) //.Background(m.background)
	if m.background {
		renderStyle = renderStyle.Background(bgColor).Foreground(fgColor)
	}
	doc.WriteString(renderStyle.Render("has style much!"))
	doc.WriteString(helpStyle.Render("\nq: exit • b: bold • s: strike • i: italic • u: underline • B: background"))

	return doc.String()
}
