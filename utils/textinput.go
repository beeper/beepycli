package utils

import (
	"github.com/charmbracelet/bubbles/textinput"
	gloss "github.com/charmbracelet/lipgloss"
)

func TextInput(placeholder string, hidden bool) textinput.Model {
	t := textinput.New()
	t.Placeholder = placeholder
	t.Prompt = ""
	t.Cursor.Style = gloss.NewStyle().Foreground(Magenta)
	t.Cursor.TextStyle = gloss.NewStyle().Foreground(Purple)
	t.TextStyle = gloss.NewStyle().Foreground(Purple)
	if hidden {
		t.EchoMode = textinput.EchoPassword
	}

	return t
}
