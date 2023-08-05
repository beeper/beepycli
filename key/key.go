package key

import (
	"fmt"
	"os"

	"github.com/figbert/beepy/utils"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbletea"
	gloss "github.com/charmbracelet/lipgloss"
)

type Model struct {
	password, file textinput.Model
	status         string
	buttonFocused  bool
}

func InitModel() Model {
	var placeholder string
	h, err := os.UserHomeDir()
	if err == nil {
		placeholder = h
	} else {
		placeholder = "/home/user/element-keys.txt"
	}

	return Model{
		password: utils.TextInput(utils.PasswordPlaceholder, true),
		file:     utils.TextInput(placeholder, false),
	}
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
			if m.buttonFocused {
				return m, utils.NextPhase
			} else if m.file.Focused() {
				m.buttonFocused = true
				m.file.Blur()
			} else if m.password.Focused() {
				m.password.Blur()
				m.file.Focus()
			} else {
				m.password.Focus()
			}
		case tea.KeyShiftTab:
			if m.buttonFocused {
				m.buttonFocused = false
				m.file.Focus()
			} else if m.file.Focused() {
				m.file.Blur()
				m.password.Focus()
			} else if m.password.Focused() {
				m.password.Blur()
			} else {
				return m, utils.PrevPhase
			}
		default:
			var cmd tea.Cmd
			if m.password.Focused() {
				m.password, cmd = m.password.Update(msg)
			} else if m.file.Focused() {
				m.file, cmd = m.file.Update(msg)
				return m, tea.Batch(cmd, prevalidate)
			}
			return m, cmd
		}
	} else if _, ok := msg.(fileLoadingMsg); ok {
		m.status = fileLoading.Render()
		return m, validate(m.file.Value())
	} else if _, ok := msg.(fileOkMsg); ok {
		m.status = fileOk.Render()
	} else if _, ok := msg.(fileErrMsg); ok {
		m.status = fileErr.Render()
	}
	return m, nil
}

func (m Model) View() string {
	italic := gloss.NewStyle().Italic(true)

	return fmt.Sprintf(
		"%s\n"+
			"You're going to need some keys for this: in Beeper Desktop, navigate\n"+
			"to %s and scroll down to the\n"+
			"%s heading. Then click the %s button.\n"+
			"The fields below will ask you for the password you gave your keys,\n"+
			"and their location on your computer.\n\n"+
			"Password: %s\n"+
			"Path to Keys: %s %s\n\n"+
			"Whole lotta cryptography gonna be happening here real soon!\n\n"+
			"%s",
		utils.Title().Render("Configure end-to-end encryption for your Beepy"),
		italic.Render("Gear > Settings > Security & Privacy"),
		italic.Render("Cryptography"),
		italic.Render("Export E2E rooms keys"),
		m.password.View(),
		m.file.View(),
		m.status,
		utils.Button(m.buttonFocused).Render("Next"),
	)
}
