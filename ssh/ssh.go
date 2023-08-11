package ssh

import (
	"fmt"

	"github.com/figbert/beepy/utils"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbletea"
)

type Model struct {
	username, host, password, confirmation textinput.Model
	buttonFocused, valid                   bool
}

func InitModel() Model {
	return Model{
		username:     utils.TextInput("beepy", false),
		host:         utils.TextInput("127.0.0.1", false),
		password:     utils.TextInput("beepbeep", true),
		confirmation: utils.TextInput("beepbeep", true),
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
				if len(m.username.Value()) > 0 && len(m.host.Value()) > 0 && m.valid {
					return m, utils.NextPhase
				}
			} else if m.confirmation.Focused() {
				m.confirmation.Blur()
				m.buttonFocused = true
			} else if m.password.Focused() {
				m.password.Blur()
				m.confirmation.Focus()
			} else if m.host.Focused() {
				m.host.Blur()
				m.password.Focus()
			} else if m.username.Focused() {
				m.username.Blur()
				m.username.CursorStart()
				m.host.Focus()
			} else {
				m.username.CursorEnd()
				m.username.Focus()
			}
		case tea.KeyShiftTab:
			if m.buttonFocused {
				m.buttonFocused = false
				m.confirmation.Focus()
			} else if m.confirmation.Focused() {
				m.confirmation.Blur()
				m.password.Focus()
			} else if m.password.Focused() {
				m.password.Blur()
				m.host.Focus()
			} else if m.host.Focused() {
				m.username.CursorEnd()
				m.host.Blur()
				m.username.Focus()
			} else if m.username.Focused() {
				m.username.Blur()
				m.username.CursorStart()
			} else {
				return m, utils.PrevPhase
			}
		default:
			var cmd tea.Cmd
			if m.username.Focused() {
				m.username, cmd = m.username.Update(msg)
			} else if m.host.Focused() {
				m.host, cmd = m.host.Update(msg)
			} else if m.password.Focused() {
				m.password, cmd = m.password.Update(msg)
				cmd = tea.Batch(cmd, validatePassword(m.password.Value(), m.confirmation.Value()))
			} else if m.confirmation.Focused() {
				m.confirmation, cmd = m.confirmation.Update(msg)
				cmd = tea.Batch(cmd, validatePassword(m.password.Value(), m.confirmation.Value()))
			}
			return m, cmd
		}
	} else if validation, ok := msg.(passwordValidationMsg); ok {
		m.valid = bool(validation)
	}
	return m, nil
}

func (m Model) View() string {
	validation := ""
	if len(m.confirmation.Value()) > 0 {
		if m.valid {
			validation = fmt.Sprintf(" %s", okValidation.Render())
		} else {
			validation = fmt.Sprintf(" %s", errValidation.Render())
		}
	}

	return fmt.Sprintf(
		"%s\n"+
			"The Wizard accesses your device over SSH.\n\n"+
			"\tssh %s@%s\n\n"+
			"\tPassword: %s%s\n"+
			"\tConfirm Password: %s%s\n\n"+
			"He also loves casting spells 🪄.\n\n"+
			"%s",
		utils.Title().Render("Configure network access to your Beepy"),
		m.username.View(), m.host.View(),
		m.password.View(), validation,
		m.confirmation.View(), validation,
		utils.Button(m.buttonFocused).Render("Next"),
	)
}
