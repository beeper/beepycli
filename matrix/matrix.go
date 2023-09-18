package matrix

import (
	"fmt"

	"github.com/beeper/beepycli/utils"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbletea"
)

type Model struct {
	email, code   textinput.Model
	session       string
	buttonFocused bool
}

func InitModel() Model {
	return Model{
		email: utils.TextInput("example@example.com", false),
		code:  utils.TextInput("123456", false),
	}
}

func (m Model) Session() string {
	return m.session
}

func (m Model) Code() string {
	return m.code.Value()
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
				if len(m.email.Value()) > 0 && len(m.code.Value()) > 0 {
					return m, utils.NextPhase
				}
			} else if m.code.Focused() {
				m.buttonFocused = true
				m.code.Blur()
			} else if m.email.Focused() {
				m.email.Blur()
				m.code.Focus()
				if len(m.code.Value()) == 0 {
					return m, initAuth()
				}
			} else {
				m.email.Focus()
			}
		case tea.KeyShiftTab:
			if m.buttonFocused {
				m.buttonFocused = false
				m.code.Focus()
			} else if m.code.Focused() {
				m.code.Blur()
				m.email.Focus()
			} else if m.email.Focused() {
				m.email.Blur()
			} else {
				return m, utils.PrevPhase
			}
		default:
			var cmd tea.Cmd
			if m.code.Focused() {
				switch key.String() {
				case "1", "2", "3", "4", "5", "6", "7", "8", "9", "0", "backspace":
					m.code, cmd = m.code.Update(msg)
				}
			} else if m.email.Focused() {
				m.email, cmd = m.email.Update(msg)
			}
			return m, cmd
		}
	} else if session, ok := msg.(loginStarted); ok {
		m.session = string(session)
		m.code.Placeholder = "Talking to Beeper servers…"
		return m, sendEmail(m.session, m.email.Value())
	} else if _, ok := msg.(emailSuccess); ok {
		m.code.Placeholder = "Check your inbox…"
	} else if err, ok := msg.(apiError); ok {
		m.code.Placeholder = err.Error()
	}
	return m, nil
}

func (m Model) View() string {
	return fmt.Sprintf(
		"%s\n"+
			"Let's bootstrap your Beepy... with Beeper!\n\n"+
			"\tEmail: %s\n"+
			"\tConfirmation Code: %s\n\n"+
			"We'll have you up and chatting in style in no time at all 💬.\n\n"+
			"%s",
		utils.Title().Render("Configure your Matrix account"),
		m.email.View(),
		m.code.View(),
		utils.Button(m.buttonFocused).Render("Next"),
	)
}
