package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/SpokaneTech/go-meetup-tui/pkg/db"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea" // obnoxious!
	"github.com/charmbracelet/lipgloss"
)

func main() {
	m := initModel()

	_, err := tea.NewProgram(m).Run()
	if err != nil {
		log.Fatalf("error with program: %s", err.Error())
	}
}

func initModel() Model {
	ta := textarea.New()
	ta.Placeholder = "enter query here..."
	ta.Prompt = "┃ "
	ta.CharLimit = 500
	ta.SetWidth(50)
	ta.SetHeight(5)
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()
	ta.ShowLineNumbers = false
	ta.KeyMap.InsertNewline.SetEnabled(false)

	vp := viewport.New(50, 5)

	ti := textinput.New()
	ti.Placeholder = "test.db"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return Model{
		open:     ti,
		viewport: vp,
		query:    ta,
	}
}

// Model is the main container for our tui.
// It must satisfy the tea `model` interface
type Model struct {
	viewport viewport.Model
	query    textarea.Model
	open     textinput.Model

	err string

	db db.DB
}

// Init the tui. The returned tea.Cmd represents any initial
// I/O or setup type functions we want to execute
func (m Model) Init() tea.Cmd {
	return nil
}

// Update the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var oCmd, dbCmd, queryCmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()
		switch key {
		case "ctrl+c", "q":
			// return the model, and ask for tea to quit
			return m, tea.Quit

		case "esc":
			m.query.SetValue("")

		case "enter":
			m.err = ""
			if m.db == nil {
				m.open.Blur()
				dbCmd = openDB(m.open.Value())
			} else {
				m.viewport.SetContent("")
				m.query.Blur()
				dbCmd = queryDB(m.db, m.query.Value())
			}
		}

	case DBOpen:
		m.db = msg.db
		m.open.Blur()
		m.query.Focus()

	case DBQuery:
		m.query.Focus()
		d, _ := json.Marshal(msg.rows)
		m.viewport.SetContent(string(d))

	case errMsg:
		if m.db == nil {
			m.open.Focus()
		} else {
			m.query.Focus()
			m.err = msg.err
		}
	}

	m.open, oCmd = m.open.Update(msg)
	m.query, queryCmd = m.query.Update(msg)

	return m, tea.Batch(oCmd, queryCmd, dbCmd)
}

// View displays the model
func (m Model) View() string {
	doc := strings.Builder{}

	if m.db == nil {
		doc.WriteString(fmt.Sprintf("Please enter a DB filename to open:\n\n%s\n\n(q to quit)", m.open.View()))
	} else {
		doc.WriteString(fmt.Sprintf("\n%s\n\n%s\n\n", m.query.View(), m.viewport.View()))
		if m.err != "" {
			doc.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Render(fmt.Sprintf("error: %s\n", m.err)))
		}
	}

	return doc.String()
}

////
// Our Custom Logic lives here
///

// errMsg for handling error returns
type errMsg struct {
	err string
}

// For opening our DB
type DBOpen struct {
	db db.DB
}

func openDB(path string) func() tea.Msg {
	return func() tea.Msg {
		db, err := db.New(fmt.Sprintf("sqlite:%s", path))
		if err != nil {
			return errMsg{err.Error()}
		}

		err = db.Open()
		if err != nil {
			return errMsg{err.Error()}
		}

		return DBOpen{db: db}
	}
}

// For querying our DB
type DBQuery struct {
	colnames []string
	rows     []interface{}
}

func queryDB(db db.DB, query string) func() tea.Msg {
	return func() tea.Msg {
		c, r, err := db.Query(query)
		if err != nil {
			return errMsg{err.Error()}
		}

		return DBQuery{colnames: c, rows: r}
	}
}
