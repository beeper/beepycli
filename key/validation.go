package key

import (
	"os"

	"github.com/figbert/beepy/utils"

	"github.com/charmbracelet/bubbletea"
	gloss "github.com/charmbracelet/lipgloss"
)

var (
	okValidation = gloss.NewStyle().
			Foreground(utils.Green).
			SetString("✔")
	loadingValidation = gloss.NewStyle().
				Foreground(utils.Yellow).
				SetString("↻")
	errValidation = gloss.NewStyle().
			Foreground(utils.Red).
			SetString("✘")
)

type fileOkMsg struct{}
type fileLoadingMsg struct{}
type fileErrMsg struct{}

func prevalidate() tea.Msg {
	return fileLoadingMsg{}
}

func validateKey(path string) tea.Cmd {
	return func() tea.Msg {
		file, err := os.Open(path)
		if err != nil {
			return fileErrMsg{}
		}

		info, err := file.Stat()
		if err != nil || info.IsDir() {
			return fileErrMsg{}
		}

		return fileOkMsg{}
	}
}

type passwordValidationMsg bool

func validatePassword(password, confirmation string) tea.Cmd {
	return func() tea.Msg {
		return passwordValidationMsg(len(password) > 0 && password == confirmation)
	}
}
