package key

import (
	"os"

	"github.com/figbert/beepy/utils"

	"github.com/charmbracelet/bubbletea"
	gloss "github.com/charmbracelet/lipgloss"
)

var (
	fileOk = gloss.NewStyle().
		Foreground(utils.Green).
		SetString("✔")
	fileLoading = gloss.NewStyle().
			Foreground(utils.Yellow).
			SetString("↻")
	fileErr = gloss.NewStyle().
		Foreground(utils.Red).
		SetString("✘")
)

type fileOkMsg struct{}
type fileLoadingMsg struct{}
type fileErrMsg struct{}

func prevalidate() tea.Msg {
	return fileLoadingMsg{}
}

func validate(path string) tea.Cmd {
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
