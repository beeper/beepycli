package matrix

import (
	"github.com/figbert/beepy/utils"

	"github.com/charmbracelet/bubbletea"
	gloss "github.com/charmbracelet/lipgloss"
)

var (
	passwordsValid = gloss.NewStyle().
			Foreground(utils.Green).
			SetString("✔")
	passwordsInvalid = gloss.NewStyle().
				Foreground(utils.Red).
				SetString("✘")
)

type passwordValidationMsg bool

func validate(password, confirmation string) tea.Cmd {
	return func() tea.Msg {
		return passwordValidationMsg(len(password) > 0 && password == confirmation)
	}
}
