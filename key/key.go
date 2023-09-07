package key

import (
	"fmt"

	"github.com/beeper/beepycli/utils"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbletea"
	gloss "github.com/charmbracelet/lipgloss"
)

type valid struct {
	file, password bool
}

type Model struct {
	password, confirmation, file textinput.Model
	status                       string
	valid                        valid
	buttonFocused                bool
}

func InitModel() Model {
	return Model{
		password:     utils.TextInput(utils.PasswordPlaceholder, true),
		confirmation: utils.TextInput(utils.PasswordPlaceholder, true),
		file:         utils.TextInput("/home/user/element-keys.txt", false),
	}
}

func (m Model) KeyPath() string {
	return m.file.Value()
}

func (m Model) KeyPassword() string {
	return m.password.Value()
}

func (m Model) Init() tea.Cmd {
	return getTextFileInWd
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyEnter, tea.KeyTab:
			if m.buttonFocused {
				if m.valid.file && m.valid.password {
					return m, utils.NextPhase
				}
			} else if m.file.Focused() {
				m.buttonFocused = true
				m.file.Blur()
			} else if m.confirmation.Focused() {
				m.confirmation.Blur()
				m.file.Focus()
			} else if m.password.Focused() {
				m.password.Blur()
				m.confirmation.Focus()
			} else {
				m.password.Focus()
			}
		case tea.KeyShiftTab:
			if m.buttonFocused {
				m.buttonFocused = false
				m.file.Focus()
			} else if m.file.Focused() {
				m.file.Blur()
				m.confirmation.Focus()
			} else if m.confirmation.Focused() {
				m.confirmation.Blur()
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
				cmd = tea.Batch(cmd, validatePassword(m.password.Value(), m.confirmation.Value()))
			} else if m.confirmation.Focused() {
				m.confirmation, cmd = m.confirmation.Update(msg)
				cmd = tea.Batch(cmd, validatePassword(m.password.Value(), m.confirmation.Value()))
			} else if m.file.Focused() {
				m.file, cmd = m.file.Update(msg)
				return m, tea.Batch(cmd, prevalidate)
			}
			return m, cmd
		}
	} else if _, ok := msg.(fileLoadingMsg); ok {
		m.status = loadingValidation.Render()
		m.valid.file = false
		return m, validateKey(m.file.Value())
	} else if _, ok := msg.(fileOkMsg); ok {
		m.status = okValidation.Render()
		m.valid.file = true
	} else if _, ok := msg.(fileErrMsg); ok {
		m.status = errValidation.Render()
		m.valid.file = false
	} else if validation, ok := msg.(passwordValidationMsg); ok {
		m.valid.password = bool(validation)
	} else if file, ok := msg.(keyFileMsg); ok {
		m.file.Placeholder = string(file)
		m.file.SetValue(string(file))
		m.file.CursorEnd()
		return m, prevalidate
	}
	return m, nil
}

func (m Model) View() string {
	italic := gloss.NewStyle().Italic(true)

	validation := ""
	if len(m.confirmation.Value()) > 0 {
		if m.valid.password {
			validation = fmt.Sprintf(" %s", okValidation.Render())
		} else {
			validation = fmt.Sprintf(" %s", errValidation.Render())
		}
	}

	return fmt.Sprintf(
		"%s\n"+
			"You're going to need some keys for this: in Beeper Desktop, navigate\n"+
			"to %s and scroll down to the\n"+
			"%s heading. Then click the %s button.\n"+
			"The fields below will ask you for the password you gave your keys,\n"+
			"and their location on your computer.\n\n"+
			"Password: %s%s\n"+
			"Confirm Password: %s%s\n"+
			"Path to Keys: %s %s\n\n"+
			"Whole lotta cryptography gonna be happening here real soon 🔐!\n\n"+
			"%s",
		utils.Title().Render("Configure end-to-end encryption for your Beepy"),
		italic.Render("Gear > Settings > Security & Privacy"),
		italic.Render("Cryptography"),
		italic.Render("Export E2E room keys"),
		m.password.View(), validation,
		m.confirmation.View(), validation,
		m.file.View(),
		m.status,
		utils.Button(m.buttonFocused).Render("Next"),
	)
}
