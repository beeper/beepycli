package main

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbletea"
	gloss "github.com/charmbracelet/lipgloss"
)

func configureTextInput(t *textinput.Model, placeholder string, hidden bool) {
	t.Placeholder = placeholder
	t.Prompt = ""
	t.Cursor.Style = gloss.NewStyle().Foreground(magenta)
	t.Cursor.TextStyle = gloss.NewStyle().Foreground(purple)
	t.TextStyle = gloss.NewStyle().Foreground(purple)
	if hidden {
		t.EchoMode = textinput.EchoPassword
	}
}

type sshModel struct {
	username, host textinput.Model
}

func initSSHModel() sshModel {
	m := sshModel{
		username: textinput.New(),
		host:     textinput.New(),
	}

	configureTextInput(&m.username, "user", false)
	configureTextInput(&m.host, "127.0.0.1", false)

	return m
}

type matrixModel struct {
	username, password, homeserver textinput.Model
}

func initMatrixModel() matrixModel {
	m := matrixModel{
		username:   textinput.New(),
		password:   textinput.New(),
		homeserver: textinput.New(),
	}

	configureTextInput(&m.username, "@user:example.com", false)
	configureTextInput(&m.password, passwordPlaceholder, true)
	configureTextInput(&m.homeserver, domainPlaceholder, false)

	return m
}

type model struct {
	phase         phase
	focusOnButton bool
	ssh           sshModel
	matrix        matrixModel
}

func initModel() model {
	m := model{
		phase:  welcome,
		ssh:    initSSHModel(),
		matrix: initMatrixModel(),
	}

	return m
}

func (m model) Init() tea.Cmd {
	return nil
}
