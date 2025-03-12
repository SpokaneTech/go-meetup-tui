package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/SpokaneTech/go-meetup-tui/pkg/db"
	"github.com/charmbracelet/bubbles/table"
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

	rt := table.New()
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false).
		Align(lipgloss.Left)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	rt.SetStyles(s)

	return Model{
		results:  rt,
		open:     ti,
		viewport: vp,
		query:    ta,
	}
}

// Model is the main container for our tui.
// It must satisfy the tea `model` interface
type Model struct {
	viewport viewport.Model
	results  table.Model
	query    textarea.Model
	open     textinput.Model

	err string
	msg string

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
		m.err = ""
		key := msg.String()
		switch key {
		case "ctrl+c":
			// return the model, and ask for tea to quit
			return m, tea.Quit

		case "esc":
			m.query.SetValue("")

		case "enter":
			m.msg = ""
			if m.db == nil {
				m.open.Blur()
				dbCmd = openDB(m.open.Value())
			} else {
				m.query.Blur()
				dbCmd = queryDB(m.db, m.query.Value())
			}
		}

	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width
		m.query.SetWidth(msg.Width)
		m.viewport.Height = (msg.Height / 2) - m.query.Height()
		m.results.SetHeight((msg.Height / 2) - m.query.Height())

	case DBOpen:
		m.db = msg.db
		m.open.Blur()
		m.query.Focus()

	case DBQuery:
		m.query.Focus()
		if len(msg.rows) == 0 {
			m.msg = "success"
		} else {
			cols := []table.Column{}
			for _, c := range msg.colnames {
				cols = append(cols, table.Column{Title: c, Width: lipgloss.Width(c) + 6})
			}

			rows := []table.Row{}
			for _, r := range msg.rows {
				rows = append(rows, convertRow(r, msg.colnames))
			}
			m.results.SetColumns(cols)
			m.results.SetRows(rows)
		}

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
		doc.WriteString(fmt.Sprintf("Please enter a DB filename to open:\n\n%s\n\n(ctrl+c to quit)", m.open.View()))
	} else {
		doc.WriteString(fmt.Sprintf("\n%s\n\n", m.query.View()))

		if m.err != "" {
			doc.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Render(fmt.Sprintf("error: %s\n", m.err)))
		}

		if m.msg != "" {
			doc.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#04B575")).Render(fmt.Sprintf("%s\n", m.msg)))
		} else {
			doc.WriteString(fmt.Sprintf("\n\n%s\n\n", m.results.View()))
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

func convertRow(r interface{}, cnames []string) table.Row {
	mr, ok := r.(map[string]interface{})
	if !ok {
		return nil
	}

	row := []string{}
	for _, n := range cnames {
		switch v := mr[n].(type) {
		case bool:
			row = append(row, strconv.FormatBool(v))
		case string:
			row = append(row, v)
		case int64:
			row = append(row, strconv.FormatInt(v, 0))
		case int32:
			row = append(row, string(v))
		case float64:
			row = append(row, strconv.FormatFloat(v, 'f', -1, 64))
		}
	}

	return row
}
