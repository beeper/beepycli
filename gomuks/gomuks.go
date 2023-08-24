package gomuks

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/figbert/beepy/utils"

	"maunium.net/go/gomuks/headless"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbletea"
	gloss "github.com/charmbracelet/lipgloss"
)

type mode int

const (
	loading mode = iota
	success
	failure
)

var (
	defaultMsg = utils.Title().Render("🔮 The Wizard must now gaze into his crystal ball 🔮")
)

type Model struct {
	config  headless.Config
	updates chan fmt.Stringer

	status mode

	spinner spinner.Model
	msg     string

	err error
}

func InitModel() Model {
	return Model{
		updates: make(chan fmt.Stringer),
		config:  headless.Config{},
		spinner: spinner.New(
			spinner.WithSpinner(spinner.Dot),
			spinner.WithStyle(gloss.NewStyle().Foreground(utils.Magenta)),
		),
		msg:    defaultMsg,
		status: loading,
	}
}

type loadingMsg string

func awaitLoadingMsg(updates chan fmt.Stringer) tea.Cmd {
	return func() tea.Msg {
		msg, ok := <-updates
		if ok {
			return loadingMsg(msg.String())
		} else {
			return nil
		}
	}
}

type successMsg struct{}
type failureMsg error

func initializeGomuksInstance(conf headless.Config, updates chan fmt.Stringer) tea.Cmd {
	return func() tea.Msg {
		err := headless.Init(conf, updates)
		if err != nil {
			return failureMsg(err)
		}

		return successMsg{}
	}
}

func (m Model) UpdateConfig(session, code, keyPath, keyPassword, recoveryCode string) Model {
	m.updates = make(chan fmt.Stringer)
	m.config.OutputDir = filepath.Join(os.TempDir(), "beepy", fmt.Sprintf("%d", time.Now().Unix()))
	m.config.Session = session
	m.config.Code = code
	m.config.KeyPath = keyPath
	m.config.KeyPassword = keyPassword
	m.config.RecoveryCode = recoveryCode

	return m
}

func (m Model) OutputDir() string {
	return m.config.OutputDir
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		awaitLoadingMsg(m.updates),
		initializeGomuksInstance(m.config, m.updates),
		m.spinner.Tick,
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok {
		if key.Type == tea.KeyCtrlC {
			return m, tea.Quit
		} else if m.status == failure {
			return m, utils.PrevPhase
		} else if m.status == success {
			return m, utils.NextPhase
		}
	} else if _, ok := msg.(spinner.TickMsg); ok {
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	} else if update, ok := msg.(loadingMsg); ok {
		m.msg = string(update)
		return m, awaitLoadingMsg(m.updates)
	} else if _, ok := msg.(successMsg); ok {
		m.msg = defaultMsg
		m.status = success
	} else if err, ok := msg.(failureMsg); ok {
		m.msg = defaultMsg
		m.err = err
		m.status = failure
	}
	return m, nil
}

func (m Model) View() string {
	switch m.status {
	case loading:
		return fmt.Sprintf("%s%s %s", m.spinner.View(), m.msg, m.spinner.View())
	case failure:
		return fmt.Sprintf(
			"%s\n"+
				"The wizard has divined the following: %s\n"+
				"Press any key to try again.",
			utils.Title().Render("Gomuks initialization failed"),
			utils.Error(m.err.Error()),
		)
	case success:
		return fmt.Sprintf(
			"%s\n"+
				"Well done, apprentice.\n"+
				"Press any key to continue.",
			utils.Title().Render("Gomuks initialization successful"),
		)
	}

	return ""
}
