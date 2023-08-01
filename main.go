package main

import (
	"strings"

	"github.com/charmbracelet/bubbletea"
	gloss "github.com/charmbracelet/lipgloss"

	"github.com/charmbracelet/bubbles/textinput"
)

type phase int

const (
	welcome phase = iota
	ssh
	matrix
	next
)

type sshModel struct {
	username, host textinput.Model
}

type matrixModel struct {
	username, password, homeserver textinput.Model
}

type model struct {
	phase         phase
	focusOnButton bool
	ssh           sshModel
	matrix        matrixModel
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyEnter, tea.KeyTab:
			switch m.phase {
			case welcome:
				if !m.focusOnButton {
					m.focusOnButton = true
				} else {
					m.phase = ssh
					m.focusOnButton = false
				}
			case ssh:
				if m.ssh.host.Focused() {
					m.focusOnButton = true
					m.ssh.host.Blur()
				} else if m.ssh.username.Focused() {
					m.ssh.username.Blur()
					m.ssh.username.CursorStart()
					m.ssh.host.Focus()
				} else if m.focusOnButton {
					m.phase = matrix
					m.focusOnButton = false
				} else {
					m.ssh.username.CursorEnd()
					m.ssh.username.Focus()
				}
			case matrix:
				if m.matrix.homeserver.Focused() {
					m.focusOnButton = true
					m.matrix.homeserver.Blur()
				} else if m.matrix.password.Focused() {
					m.matrix.password.Blur()
					m.matrix.homeserver.Focus()
				} else if m.matrix.username.Focused() {
					m.matrix.username.Blur()
					m.matrix.password.Focus()
				} else if m.focusOnButton {
					m.phase = next
					m.focusOnButton = false
				} else {
					m.matrix.username.Focus()
				}
			}
		case tea.KeyShiftTab:
			switch m.phase {
			case welcome:
				m.focusOnButton = false
			case ssh:
				if m.focusOnButton {
					m.focusOnButton = false
					m.ssh.host.Focus()
				} else if m.ssh.host.Focused() {
					m.ssh.username.CursorEnd()
					m.ssh.host.Blur()
					m.ssh.username.Focus()
				} else if m.ssh.username.Focused() {
					m.ssh.username.Blur()
					m.ssh.username.CursorStart()
				} else {
					m.phase = welcome
					m.focusOnButton = true
				}
			case matrix:
				if m.focusOnButton {
					m.focusOnButton = false
					m.matrix.homeserver.Focus()
				} else if m.matrix.homeserver.Focused() {
					m.matrix.homeserver.Blur()
					m.matrix.password.Focus()
				} else if m.matrix.password.Focused() {
					m.matrix.password.Blur()
					m.matrix.username.Focus()
				} else if m.matrix.username.Focused() {
					m.matrix.username.Blur()
				} else {
					m.phase = ssh
					m.focusOnButton = true
				}
			}
		default:
			var cmd tea.Cmd
			if m.phase == ssh {
				if m.ssh.username.Focused() {
					m.ssh.username, cmd = m.ssh.username.Update(msg)
				} else if m.ssh.host.Focused() {
					m.ssh.host, cmd = m.ssh.host.Update(msg)
				}
			} else if m.phase == matrix {
				if m.matrix.homeserver.Focused() {
					m.matrix.homeserver, cmd = m.matrix.homeserver.Update(msg)
				} else if m.matrix.password.Focused() {
					m.matrix.password, cmd = m.matrix.password.Update(msg)
				} else if m.matrix.username.Focused() {
					m.matrix.username, cmd = m.matrix.username.Update(msg)
				}
			}
			return m, cmd
		}
	}
	return m, nil
}

func (m model) View() string {
	var builder strings.Builder

	title := gloss.NewStyle().Bold(true)
	button := gloss.NewStyle().
		Background(gloss.Color("8")).
		Padding(0, 2)
	if m.focusOnButton {
		button = button.Background(gloss.Color("13"))
	}

	switch m.phase {
	case welcome:
		builder.WriteString(title.Render("Hello! Welcome to the Beepy Setup Wizard™"))
		builder.WriteString("\nA quick guide to navigating the Wizard:\n\n")
		builder.WriteString("\t↹ Tab|⏎ Return\n")
		builder.WriteString("\t\tMove focus forward, or progress to the next page\n")
		builder.WriteString("\t⇧↹ Shift-Tab\n")
		builder.WriteString("\t\tMove focus backward, or return to the previous page\n")
		builder.WriteString("\t^C Ctrl-C\n")
		builder.WriteString("\t\tQuit\n\n")
		builder.WriteString("We hope you enjoy your time with the Wizard 🧙!\n\n")
		builder.WriteString(button.Render("Next"))
	case ssh:
		builder.WriteString(title.Render("Configure network access to your Beepy"))
		builder.WriteString("\nThe Wizard accesses your device over SSH.\n\n")
		builder.WriteString("\tssh " + m.ssh.username.View() + "@" + m.ssh.host.View() + "\n\n")
		builder.WriteString("This is extremely required.\n\n")
		builder.WriteString(button.Render("Next"))
	case matrix:
		builder.WriteString(title.Render("Configure your Matrix account"))
		builder.WriteString("\nLet's bootstrap your Beepy... with Beeper!\n\n")
		builder.WriteString("\tUsername: " + m.matrix.username.View() + "\n\n")
		builder.WriteString("\tPassword: " + m.matrix.password.View() + "\n\n")
		builder.WriteString("\tHomeserver: " + m.matrix.homeserver.View() + "\n\n")
		builder.WriteString("We'll have you up and chatting in style in no time at all.\n\n")
		builder.WriteString(button.Render("Next"))
	default:
		builder.WriteString("How did we get here?")
	}

	return builder.String()
}

func main() {
	m := model{
		phase: welcome,
		ssh: sshModel{
			username: textinput.New(),
			host:     textinput.New(),
		},
		matrix: matrixModel{
			username:   textinput.New(),
			password:   textinput.New(),
			homeserver: textinput.New(),
		},
	}

	m.ssh.username.Placeholder = "user"
	m.ssh.username.Prompt = ""
	m.ssh.username.Cursor.Style = gloss.NewStyle().Foreground(gloss.Color("13"))
	m.ssh.host.Placeholder = "192.0.0.1"
	m.ssh.host.Prompt = ""
	m.ssh.username.Cursor.Style = gloss.NewStyle().Foreground(gloss.Color("13"))

	m.matrix.username.Placeholder = "@user:example.com"
	m.matrix.username.Prompt = ""
	m.matrix.username.Cursor.Style = gloss.NewStyle().Foreground(gloss.Color("13"))
	m.matrix.password.Placeholder = "correct horse battery staple"
	m.matrix.password.Prompt = ""
	m.matrix.password.Cursor.Style = gloss.NewStyle().Foreground(gloss.Color("13"))
	m.matrix.password.EchoMode = textinput.EchoPassword
	m.matrix.homeserver.Placeholder = "https://example.com"
	m.matrix.homeserver.Prompt = ""
	m.matrix.homeserver.Cursor.Style = gloss.NewStyle().Foreground(gloss.Color("13"))

	tea.NewProgram(m).Run()
}
