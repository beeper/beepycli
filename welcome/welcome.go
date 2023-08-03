package welcome

import (
	"fmt"

	"github.com/figbert/beepy/utils"

	"github.com/charmbracelet/bubbletea"
)

type Model bool

func InitModel() Model {
	return false
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyEnter, tea.KeyTab:
			if m {
				return m, utils.NextPhase
			} else {
				m = true
			}
		case tea.KeyShiftTab:
			m = false
		}
	}
	return m, nil
}

func (m Model) View() string {
	return fmt.Sprintf(
		"%s\n"+
			"A quick guide to navigating the Wizard:\n\n"+
			"\t↹ Tab|⏎ Return\n"+
			"\t\tMove focus forward, or progress to the next page\n"+
			"\t⇧↹ Shift-Tab\n"+
			"\t\tMove focus backward, or return to the previous page\n"+
			"\t^C Ctrl-C\n"+
			"\t\tQuit\n\n"+
			"We hope you enjoy your time with the Wizard 🧙!\n\n"+
			"%s",
		utils.Title().Render("Hello! Welcome to the Beepy Setup Wizard™"),
		utils.Button(bool(m)).Render("Next"),
	)
}
