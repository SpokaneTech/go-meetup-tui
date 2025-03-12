TUIs in golang
---

We'll run through some basics to go from 0 to ...not 0 using [BubbleTea]("https://github.com/charmbracelet/bubbletea")

I've structured the project to have a bunch of separate `cmd` directories that progressively look at building a TUI. I've implemented a lot of the boilerplate but left some of the meaty logic out so we can build that together during the meetup.

The repo is intended to be explored in the following order:

* [Survey](cmd/survey)
* [Bubble Tea Basics](cmd/btea/)
* [Bubble Tea Commands](cmd/bteacmd/)
* [Lipgloss](cmd/gloss/)
* [Something Fun](cmd/sql/)