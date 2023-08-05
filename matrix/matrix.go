package matrix

import (
	"fmt"

	"github.com/figbert/beepy/utils"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbletea"
)

type Model struct {
	username, password, homeserver textinput.Model
	buttonFocused                  bool
}

func InitModel() Model {
	return Model{
		username:   utils.TextInput("@user:example.com", false),
		password:   utils.TextInput(utils.PasswordPlaceholder, true),
		homeserver: utils.TextInput(utils.DomainPlaceholder, false),
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
			} else if m.homeserver.Focused() {
				m.buttonFocused = true
				m.homeserver.Blur()
			} else if m.password.Focused() {
				m.password.Blur()
				m.homeserver.CursorEnd()
				m.homeserver.Focus()
			} else if m.username.Focused() {
				m.username.Blur()
				m.password.Focus()
				return m, getHomeserverString(m.username.Value())
			} else {
				m.username.Focus()
			}
		case tea.KeyShiftTab:
			if m.buttonFocused {
				m.buttonFocused = false
				m.homeserver.Focus()
			} else if m.homeserver.Focused() {
				m.homeserver.Blur()
				m.password.Focus()
			} else if m.password.Focused() {
				m.password.Blur()
				m.username.Focus()
			} else if m.username.Focused() {
				m.username.Blur()
			} else {
				return m, utils.PrevPhase
			}
		default:
			var cmd tea.Cmd
			if m.homeserver.Focused() {
				m.homeserver, cmd = m.homeserver.Update(msg)
			} else if m.password.Focused() {
				m.password, cmd = m.password.Update(msg)
			} else if m.username.Focused() {
				m.username, cmd = m.username.Update(msg)
			}
			return m, cmd
		}
	} else if username, ok := msg.(usernameParseMsg); ok {
		m.homeserver.Reset()
		m.homeserver.Placeholder = "Resolving..."
		return m, resolveWellKnown(string(username))
	} else if usernameErr, ok := msg.(usernameErrMsg); ok {
		m.homeserver.Reset()
		m.homeserver.Placeholder = string(usernameErr)
	} else if homeserver, ok := msg.(homeserverParseMsg); ok {
		m.homeserver.Placeholder = utils.DomainPlaceholder
		m.homeserver.SetValue(string(homeserver))
	} else if homeserverErr, ok := msg.(homeserverErrMsg); ok {
		m.homeserver.Reset()
		m.homeserver.Placeholder = string(homeserverErr)
	}
	return m, nil
}

func (m Model) View() string {
	return fmt.Sprintf(
		"%s\n"+
			"Let's bootstrap your Beepy... with Beeper!\n\n"+
			"\tUsername: %s\n\n"+
			"\tPassword: %s\n\n"+
			"\tHomeserver: %s\n\n"+
			"We'll have you up and chatting in style in no time at all 💬.\n\n"+
			"%s",
		utils.Title().Render("Configure your Matrix account"),
		m.username.View(),
		m.password.View(),
		m.homeserver.View(),
		utils.Button(m.buttonFocused).Render("Next"),
	)
}
