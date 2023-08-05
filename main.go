package main

import (
	"github.com/charmbracelet/bubbletea"

	"github.com/figbert/beepy/utils"

	"github.com/figbert/beepy/key"
	"github.com/figbert/beepy/matrix"
	"github.com/figbert/beepy/ssh"
	"github.com/figbert/beepy/verification"
	"github.com/figbert/beepy/welcome"
)

type phase int

const (
	welcomePhase phase = iota
	sshPhase
	matrixPhase
	keyPhase
	verificationPhase
)

type model struct {
	phase        phase
	welcome      welcome.Model
	ssh          ssh.Model
	matrix       matrix.Model
	key          key.Model
	verification verification.Model
}

func initModel() model {
	return model{
		phase:        welcomePhase,
		welcome:      welcome.InitModel(),
		ssh:          ssh.InitModel(),
		matrix:       matrix.InitModel(),
		key:          key.InitModel(),
		verification: verification.InitModel(),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if _, ok := msg.(utils.NextPhaseMsg); ok {
		if m.phase < verificationPhase {
			m.phase++
		}
		return m, nil
	} else if _, ok := msg.(utils.PrevPhaseMsg); ok {
		if m.phase > welcomePhase {
			m.phase--
		}
		return m, nil
	} else {
		switch m.phase {
		case welcomePhase:
			wlcm, cmd := m.welcome.Update(msg)
			m.welcome = wlcm.(welcome.Model)
			return m, cmd
		case sshPhase:
			s, cmd := m.ssh.Update(msg)
			m.ssh = s.(ssh.Model)
			return m, cmd
		case matrixPhase:
			mtrx, cmd := m.matrix.Update(msg)
			m.matrix = mtrx.(matrix.Model)
			return m, cmd
		case keyPhase:
			k, cmd := m.key.Update(msg)
			m.key = k.(key.Model)
			return m, cmd
		case verificationPhase:
			vrfy, cmd := m.verification.Update(msg)
			m.verification = vrfy.(verification.Model)
			return m, cmd
		}
	}
	return m, nil
}

func (m model) View() string {
	switch m.phase {
	case welcomePhase:
		return m.welcome.View()
	case sshPhase:
		return m.ssh.View()
	case matrixPhase:
		return m.matrix.View()
	case keyPhase:
		return m.key.View()
	case verificationPhase:
		return m.verification.View()
	default:
		return "How did we end up here?"
	}
}

func main() {
	m := initModel()
	tea.NewProgram(m).Run()
}
