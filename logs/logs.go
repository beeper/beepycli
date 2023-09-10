package logs

import (
	"fmt"

	bsh "github.com/beeper/beepycli/ssh"
	"github.com/beeper/beepycli/utils"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbletea"
	gloss "github.com/charmbracelet/lipgloss"
)

type model struct {
	exfiltrating, done bool
	message, file      string
	error              error

	setup   bsh.Model
	spinner spinner.Model

	conf   *ssh.ClientConfig
	conn   *ssh.Client
	client *sftp.Client
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case utils.NextPhaseMsg:
		m.exfiltrating = true
		return m, tea.Batch(
			genConfig(m.setup.Username(), m.setup.Password()),
			m.spinner.Tick,
		)
	case configGenerated:
		m.conf = msg
		m.message = "Generated SSH configuration…"
		return m, genConn(m.setup.Host(), m.conf)
	case connGenerated:
		m.conn = msg
		m.message = "Opened connection with Beepy…"
		return m, locateLogFileOnBeepy(m.conn)
	case logFileFound:
		m.file = string(msg)
		m.message = "Log file found…"
		return m, genClient(m.conn)
	case clientGenerated:
		m.client = msg
		m.message = "Created SFTP client…"
		return m, copyLogs(m.client, m.file)
	case logsCopiedSuccessfully:
		m.done = true
		m.message = fmt.Sprintf(
			"🪶 Copied logs from your Beepy to: %s\n🙏 May the programmer working on your bug report know peace…",
			gloss.NewStyle().Foreground(utils.Gray).Render(string(msg)),
		)
		return m, tea.Quit
	case transferErr:
		m.error = msg
		return m, tea.Quit
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	default:
		if !m.exfiltrating && !m.done {
			stp, cmd := m.setup.Update(msg)
			m.setup = stp.(bsh.Model)
			return m, cmd
		}
	}
	return m, nil
}

func (m model) View() string {
	if !m.exfiltrating {
		return m.setup.View()
	} else if m.error != nil {
		return fmt.Sprintf("❌ %s ❌\n", m.error)
	} else if !m.done {
		return fmt.Sprintf("%s%s %s\n", m.spinner.View(), m.message, m.spinner.View())
	}

	return m.message
}

func Run() {
	m := model{
		message: utils.Title().Render("🔮 The Wizard must now gaze into his crystal ball 🔮"),
		setup:   bsh.InitModel(),
		spinner: spinner.New(
			spinner.WithSpinner(spinner.Dot),
			spinner.WithStyle(gloss.NewStyle().Foreground(utils.Magenta)),
		),
	}

	tea.NewProgram(m).Run()
}
