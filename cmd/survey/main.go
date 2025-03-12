package main

import (
	"fmt"
	"log"
	"os"
	"slices"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/jedib0t/go-pretty/v6/table"
)

func main() {
	// let's ask a question
	prompt := &survey.Select{
		Message: "favorite ice cream falvor",
		Options: []string{"vanilla", "chocolate", "strawberry", "you have poor taste in ice cream (i.e. not listed)"},
	}

	var flavor string
	err := survey.AskOne(prompt, &flavor)
	if err != nil {
		if err == terminal.InterruptErr {
			fmt.Println("goodbye")
			os.Exit(0)
		}
		log.Fatal(err)
	}

	fmt.Printf("you selected: %s\n\n", flavor)

	// Output a Table
	t := table.NewWriter()
	t.SetStyle(table.StyleRounded)
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"ID", "Name"})
	for i, r := range []string{"Seneca", "Pythagoras", "Cicero", "Epictetus"} {
		t.AppendRow(table.Row{i + 1, r})
	}
	t.Render()

	// Ask Another Question
	input := &survey.Input{
		Message: "Who was your favorite philosopher? (enter their id)",
	}

	selected := ""
	for {
		err := survey.AskOne(input, &selected)
		if err != nil {
			if err == terminal.InterruptErr {
				fmt.Println("goodbye")
				os.Exit(0)
			}
			log.Fatal(err)
		}

		if slices.Contains([]string{"1", "2", "3", "4"}, selected) {
			break
		}
	}

	fmt.Println("good choice!")
}
