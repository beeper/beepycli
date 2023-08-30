package ssh

import (
	"github.com/beeper/beepycli/utils"

	"github.com/charmbracelet/bubbletea"
	gloss "github.com/charmbracelet/lipgloss"
)

var (
	okValidation = gloss.NewStyle().
			Foreground(utils.Green).
			SetString("✔")
	errValidation = gloss.NewStyle().
			Foreground(utils.Red).
			SetString("✘")
)

type passwordValidationMsg bool

func validatePassword(password, confirmation string) tea.Cmd {
	return func() tea.Msg {
		return passwordValidationMsg(len(password) > 0 && password == confirmation)
	}
}
