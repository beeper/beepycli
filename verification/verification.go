package verification

import (
	"fmt"

	"github.com/figbert/beepy/utils"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbletea"
)

type Model struct {
	recovery      textinput.Model
	buttonFocused bool
}

func InitModel() Model {
	return Model{recovery: utils.TextInput("tDAK LMRH PiYE bdzi maCe xLX5 wV6P Nmfd c5mC wLef 15Fs VVSc", true)}
}

func (m Model) RecoveryPhrase() string {
	return m.recovery.Value()
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
				if len(m.recovery.Value()) > 0 {
					return m, utils.NextPhase
				}
			} else if m.recovery.Focused() {
				m.buttonFocused = true
				m.recovery.Blur()
			} else {
				m.recovery.Focus()
			}
		case tea.KeyShiftTab:
			if m.buttonFocused {
				m.buttonFocused = false
				m.recovery.Focus()
			} else if m.recovery.Focused() {
				m.recovery.Blur()
			} else {
				return m, utils.PrevPhase
			}
		default:
			var cmd tea.Cmd
			if m.recovery.Focused() {
				m.recovery, cmd = m.recovery.Update(msg)
			}
			return m, cmd
		}
	}
	return m, nil
}

func (m Model) View() string {
	return fmt.Sprintf(
		"%s\n"+
			"Please input your account recovery passphrase. This is the 48 character\n"+
			"recovery code you received when you first set up Beeper.\n\n"+
			"Recovery Phrase: %s\n\n"+
			"Something cool about the Beepy: this step doesn't cost $8/month 🤑!\n\n"+
			"%s",
		utils.Title().Render("Verify your Beepy client"),
		m.recovery.View(),
		utils.Button(m.buttonFocused).Render("Next"),
	)
}
