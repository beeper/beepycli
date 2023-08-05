package key

import (
	"os"

	"github.com/charmbracelet/bubbletea"
	gloss "github.com/charmbracelet/lipgloss"
)

var (
	fileOk = gloss.NewStyle().
		Foreground(gloss.Color("2")).
		SetString("✔")
	fileLoading = gloss.NewStyle().
			Foreground(gloss.Color("3")).
			SetString("↻")
	fileErr = gloss.NewStyle().
		Foreground(gloss.Color("1")).
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
