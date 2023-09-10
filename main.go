package main

import (
	"flag"

	"github.com/charmbracelet/bubbletea"

	"github.com/beeper/beepycli/utils"

	"github.com/beeper/beepycli/gomuks"
	"github.com/beeper/beepycli/key"
	"github.com/beeper/beepycli/logs"
	"github.com/beeper/beepycli/matrix"
	"github.com/beeper/beepycli/ssh"
	"github.com/beeper/beepycli/transfer"
	"github.com/beeper/beepycli/verification"
	"github.com/beeper/beepycli/welcome"
)

type nowWeCanQuit struct{}

func quitStageTransfer() tea.Msg {
	return nowWeCanQuit{}
}

type phase int

const (
	welcomePhase phase = iota
	matrixPhase
	keyPhase
	verificationPhase
	gomuksPhase
	sshPhase
	transferPhase
	quitPhase
)

type model struct {
	phase phase

	welcome      welcome.Model
	matrix       matrix.Model
	key          key.Model
	verification verification.Model
	gomuks       gomuks.Model
	ssh          ssh.Model
	transfer     transfer.Model
}

func initModel() model {
	return model{
		phase:        welcomePhase,
		welcome:      welcome.InitModel(),
		matrix:       matrix.InitModel(),
		key:          key.InitModel(),
		verification: verification.InitModel(),
		gomuks:       gomuks.InitModel(),
		ssh:          ssh.InitModel(),
		transfer:     transfer.InitModel(),
	}
}

func (m model) getGomuksConfig() (string, string, string, string, string) {
	return m.matrix.Session(),
		m.matrix.Code(),
		m.key.KeyPath(),
		m.key.KeyPassword(),
		m.verification.RecoveryCode()
}

func (m model) getTransferConfig() (string, string, string, string) {
	return m.ssh.Username(),
		m.ssh.Password(),
		m.ssh.Host(),
		m.gomuks.OutputDir()
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if _, ok := msg.(utils.NextPhaseMsg); ok {
		if m.phase == sshPhase {
			m.transfer = m.transfer.UpdateConfig(m.getTransferConfig())
			m.phase++
			return m, m.transfer.Init()
		} else if m.phase == verificationPhase {
			m.gomuks = m.gomuks.UpdateConfig(m.getGomuksConfig())
			m.phase++
			return m, m.gomuks.Init()
		} else if m.phase == matrixPhase {
			m.phase++
			return m, m.key.Init()
		} else if m.phase < transferPhase {
			m.phase++
		} else {
			m.phase++
			return m, quitStageTransfer
		}
		return m, nil
	} else if _, ok := msg.(utils.PrevPhaseMsg); ok {
		if m.phase == gomuksPhase {
			m.phase = matrixPhase
		} else if m.phase > welcomePhase {
			m.phase--
		}
		return m, nil
	} else if _, ok := msg.(nowWeCanQuit); ok {
		return m, tea.Quit
	} else {
		switch m.phase {
		case welcomePhase:
			wlcm, cmd := m.welcome.Update(msg)
			m.welcome = wlcm.(welcome.Model)
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
		case gomuksPhase:
			gmks, cmd := m.gomuks.Update(msg)
			m.gomuks = gmks.(gomuks.Model)
			return m, cmd
		case sshPhase:
			s, cmd := m.ssh.Update(msg)
			m.ssh = s.(ssh.Model)
			return m, cmd
		case transferPhase:
			t, cmd := m.transfer.Update(msg)
			m.transfer = t.(transfer.Model)
			return m, cmd
		}
	}
	return m, nil
}

func (m model) View() string {
	switch m.phase {
	case welcomePhase:
		return m.welcome.View()
	case matrixPhase:
		return m.matrix.View()
	case keyPhase:
		return m.key.View()
	case verificationPhase:
		return m.verification.View()
	case gomuksPhase:
		return m.gomuks.View()
	case sshPhase:
		return m.ssh.View()
	case transferPhase:
		return m.transfer.View()
	case quitPhase:
		return "✨ Enjoy your Beepy! ✨\n   Go forth and do\n   everything else!\n"
	default:
		return "How did we end up here?"
	}
}

func main() {
	shouldRunLogExfil := flag.Bool("logs", false, "Copy logs from your Beepy device to your computer")
	flag.Parse()

	if *shouldRunLogExfil {
		logs.Run()
		return
	}

	m := initModel()
	tea.NewProgram(m).Run()
}
