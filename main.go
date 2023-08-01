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
	next
)

type model struct {
	phase          phase
	focusOnButton  bool
	username, host textinput.Model
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
				if m.host.Focused() {
					m.focusOnButton = true
					m.host.Blur()
				} else if m.username.Focused() {
					m.username.Blur()
					m.username.CursorStart()
					m.host.Focus()
				} else if m.focusOnButton {
					m.phase = next
					m.focusOnButton = false
				} else {
					m.username.Focus()
				}
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
		builder.WriteString("\t⇧↹ Shift-Tab|⇧⏎ Shift-Return\n")
		builder.WriteString("\t\tMove focus backward, or return to the previous page\n")
		builder.WriteString("\t^C Ctrl-C\n")
		builder.WriteString("\t\tQuit\n\n")
		builder.WriteString("We hope you enjoy your time with the Wizard 🧙!\n\n")
		builder.WriteString(button.Render("Next"))
	case ssh:
		builder.WriteString(title.Render("Configure network access to your Beepy"))
		builder.WriteString("\nThe Wizard accesses your device over SSH.\n\n")
		builder.WriteString("\tssh " + m.username.View() + "@" + m.host.View() + "\n\n")
		builder.WriteString("This is extremely required.\n\n")
		builder.WriteString(button.Render("Next"))
	default:
		builder.WriteString("How did we get here?")
	}

	return builder.String()
}

func main() {
	m := model{
		phase:    welcome,
		username: textinput.New(),
		host:     textinput.New(),
	}

	m.username.Placeholder = "eric"
	m.username.Prompt = ""
	m.username.Cursor.Style = gloss.NewStyle().Foreground(gloss.Color("13"))
	m.host.Placeholder = "192.0.0.1"
	m.host.Prompt = ""
	m.username.Cursor.Style = gloss.NewStyle().Foreground(gloss.Color("13"))

	tea.NewProgram(m).Run()
}
