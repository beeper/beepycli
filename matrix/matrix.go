package matrix

import (
	"fmt"

	"github.com/figbert/beepy/utils"

	"maunium.net/go/mautrix/id"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbletea"
)

type Model struct {
	localpart, homeserver,
	password, confirmation, url textinput.Model
	usernameError                 string
	buttonFocused, passwordsValid bool
}

func InitModel() Model {
	return Model{
		localpart:    utils.TextInput("user", false),
		homeserver:   utils.TextInput("example.com", false),
		password:     utils.TextInput(utils.PasswordPlaceholder, true),
		confirmation: utils.TextInput(utils.PasswordPlaceholder, true),
		url:          utils.TextInput(utils.DomainPlaceholder, false),
	}
}

func (m Model) MxID() id.UserID {
	return id.NewUserID(m.localpart.Value(), m.homeserver.Value())
}

func (m Model) MxPassword() string {
	return m.password.Value()
}

func (m Model) Homeserver() string {
	return m.url.Value()
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
				if m.usernameError == "" && m.passwordsValid && len(m.url.Value()) > 0 {
					return m, utils.NextPhase
				}
			} else if m.url.Focused() {
				m.buttonFocused = true
				m.url.Blur()
			} else if m.confirmation.Focused() {
				m.confirmation.Blur()
				m.url.CursorEnd()
				m.url.Focus()
			} else if m.password.Focused() {
				m.password.Blur()
				m.confirmation.Focus()
			} else if m.homeserver.Focused() {
				m.homeserver.Blur()
				m.password.Focus()
				return m, getHomeserverString(m.MxID())
			} else if m.localpart.Focused() {
				m.localpart.Blur()
				m.localpart.CursorStart()
				m.homeserver.Focus()
			} else {
				m.localpart.Focus()
			}
		case tea.KeyShiftTab:
			if m.buttonFocused {
				m.buttonFocused = false
				m.url.Focus()
			} else if m.url.Focused() {
				m.url.Blur()
				m.confirmation.Focus()
			} else if m.confirmation.Focused() {
				m.confirmation.Blur()
				m.password.Focus()
			} else if m.password.Focused() {
				m.password.Blur()
				m.homeserver.Focus()
			} else if m.homeserver.Focused() {
				m.homeserver.Blur()
				m.localpart.CursorEnd()
				m.localpart.Focus()
			} else if m.localpart.Focused() {
				m.localpart.Blur()
			} else {
				return m, utils.PrevPhase
			}
		default:
			var cmd tea.Cmd
			if m.url.Focused() {
				m.url, cmd = m.url.Update(msg)
			} else if m.confirmation.Focused() {
				m.confirmation, cmd = m.confirmation.Update(msg)
				cmd = tea.Batch(cmd, validate(m.password.Value(), m.confirmation.Value()))
			} else if m.password.Focused() {
				m.password, cmd = m.password.Update(msg)
				cmd = tea.Batch(cmd, validate(m.password.Value(), m.confirmation.Value()))
			} else if m.homeserver.Focused() {
				m.homeserver, cmd = m.homeserver.Update(msg)
			} else if m.localpart.Focused() {
				m.localpart, cmd = m.localpart.Update(msg)
			}
			return m, cmd
		}
	} else if username, ok := msg.(usernameParseMsg); ok {
		m.url.Reset()
		m.url.Placeholder = "Resolving..."
		m.usernameError = ""
		return m, resolveWellKnown(string(username))
	} else if usernameErr, ok := msg.(usernameErrMsg); ok {
		m.url.Reset()
		m.url.Placeholder = utils.DomainPlaceholder
		m.usernameError = string(usernameErr)
	} else if hs, ok := msg.(homeserverParseMsg); ok {
		m.url.Placeholder = utils.DomainPlaceholder
		m.url.SetValue(string(hs))
	} else if hsErr, ok := msg.(homeserverErrMsg); ok {
		m.url.Reset()
		m.url.Placeholder = string(hsErr)
	} else if validation, ok := msg.(passwordValidationMsg); ok {
		m.passwordsValid = bool(validation)
	}
	return m, nil
}

func (m Model) View() string {
	err := "\n"
	if m.usernameError != "" {
		err = "\t" + utils.Error(m.usernameError) + "\n\n"
	}

	validation := ""
	if len(m.confirmation.Value()) > 0 {
		if m.passwordsValid {
			validation = fmt.Sprintf(" %s", passwordsValid.Render())
		} else {
			validation = fmt.Sprintf(" %s", passwordsInvalid.Render())
		}
	}

	return fmt.Sprintf(
		"%s\n"+
			"Let's bootstrap your Beepy... with Beeper!\n\n"+
			"\tUsername: @%s:%s\n%s"+
			"\tPassword: %s%s\n"+
			"\tConfirm Password: %s%s\n\n"+
			"\tHomeserver: %s\n\n"+
			"We'll have you up and chatting in style in no time at all 💬.\n\n"+
			"%s",
		utils.Title().Render("Configure your Matrix account"),
		m.localpart.View(),
		m.homeserver.View(),
		err,
		m.password.View(),
		validation,
		m.confirmation.View(),
		validation,
		m.url.View(),
		utils.Button(m.buttonFocused).Render("Next"),
	)
}
