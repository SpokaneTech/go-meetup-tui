package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	tea "github.com/charmbracelet/bubbletea" // obnoxious!

	"github.com/tidwall/gjson"
)

const (
	weather = "https://api.weather.gov/gridpoints/OTX/140,90/forecast"
)

func main() {
	m := &Model{message: "checking the weather..."}

	_, err := tea.NewProgram(m).Run()
	if err != nil {
		log.Fatalf("error with program: %s", err.Error())
	}
}

// Model is the main container for our tui.
// It must satisfy the tea `model` interface
type Model struct {
	message  string
	forecast string
}

// Init the tui. The returned tea.Cmd represents any initial
// I/O or setup type functions we want to execute
func (m Model) Init() tea.Cmd {
	return getForecast
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
		case "f5":
			// let's trigger and update
			m.message = "refreshing"
			return m, getForecast
		}
	case Forecast:
		if msg.Error != nil {
			m.message = fmt.Sprintf("got an error: %s", msg.Error.Error())
		} else {
			m.forecast = msg.Data
			m.message = ""
		}
	}

	return m, nil
}

func (m Model) View() string {
	if m.message != "" {
		return m.message
	}

	return fmt.Sprintf("Your Forecast For Today (f5 to refresh):\n\n%s", m.forecast)
}

// Forecast is what we are going to return
// as our `tea.Msg`. We use a type here
// so we can understand what is being
// sent in the update
type Forecast struct {
	Error error
	Data  string
}

// getForecast is our BubbleTea Command
// Something we want to do `async` and then
// process the results as an update
func getForecast() tea.Msg {
	c := &http.Client{Timeout: 10 * time.Second}
	res, err := c.Get(weather)
	if err != nil {
		return Forecast{Error: err}
	}

	bdy, err := io.ReadAll(res.Body)
	if err != nil {
		return Forecast{Error: err}
	}

	periods := gjson.Get(string(bdy), "properties.periods.0.detailedForecast")
	return Forecast{Data: periods.Raw}
}
