package ssh

import (
	"fmt"

	"github.com/figbert/beepy/utils"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbletea"
)

type Model struct {
	username, host textinput.Model
	buttonFocused  bool
}

func InitModel() Model {
	return Model{
		username: utils.TextInput("user", false),
		host:     utils.TextInput("127.0.0.1", false),
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
			} else if m.host.Focused() {
				m.buttonFocused = true
				m.host.Blur()
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
			}
			return m, cmd
		}
	}
	return m, nil
}

func (m Model) View() string {
	return fmt.Sprintf(
		"%s\n"+
			"The Wizard accesses your device over SSH.\n\n"+
			"\tssh %s@%s\n\n"+
			"He also loves casting spells 🪄.\n\n"+
			"%s",
		utils.Title().Render("Configure network access to your Beepy"),
		m.username.View(),
		m.host.View(),
		utils.Button(m.buttonFocused).Render("Next"),
	)
}
