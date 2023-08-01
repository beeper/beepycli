package main

import (
	"strings"

	"github.com/charmbracelet/bubbletea"
	gloss "github.com/charmbracelet/lipgloss"
)

type phase int

const (
	welcome phase = iota
	ssh
	matrix
	next

	passwordPlaceholder = "correct horse battery staple"
	domainPlaceholder   = "https://example.com"
)

var (
	magenta = gloss.Color("13")
	purple  = gloss.Color("5")
)

func (m model) Init() tea.Cmd {
	return nil
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
	m := initModel()
	tea.NewProgram(m).Run()
}
