package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea" // obnoxious!
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
type Model struct{}

// Init the tui. The returned tea.Cmd represents any initial
// I/O or setup type functions we want to execute
func (m Model) Init() tea.Cmd {
	return nil
}

// Update the model. This is called whenever a cmd
// in the app is run and a `tea.Msg` is created
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()
		switch key {
		case "ctrl+c", "q":
			// return the model, and ask for tea to quit
			return m, tea.Quit
		default:
			// let's set a message to display

		}
	}

	return m, nil
}

func (m Model) View() string {
	// Let's return the message to display!
	return ""
}
