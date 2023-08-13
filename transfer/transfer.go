package transfer

import (
	"fmt"

	"github.com/figbert/beepy/utils"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"

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

type transferErr error

type collection struct {
	conf, data, cache, state, media string
}

type Model struct {
	username, password, host, root string

	spinner spinner.Model

	conf   *ssh.ClientConfig
	conn   *ssh.Client
	client *sftp.Client

	dirs, patterns *collection

	status mode
	msg    string
	err    error
}

func InitModel() Model {
	return Model{
		spinner: spinner.New(
			spinner.WithSpinner(spinner.Dot),
			spinner.WithStyle(gloss.NewStyle().Foreground(utils.Magenta)),
		),
		msg:    defaultMsg,
		status: loading,
	}
}

func (m Model) UpdateConfig(username, password, host, root string) Model {
	m.username = username
	m.password = password
	m.host = host
	m.root = root

	return m
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		genConfig(m.username, m.password),
		genDirCollection(m.username),
		genPatternCollection(m.root),
		m.spinner.Tick,
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if conf, ok := msg.(configGenerated); ok {
		m.conf = conf
		m.msg = "Generated SSH configuration…"
		return m, genConn(m.host, m.conf)
	} else if dirs, ok := msg.(dirCollectionGenerated); ok {
		m.dirs = dirs
	} else if patterns, ok := msg.(patternCollectionGenerated); ok {
		m.patterns = patterns
	} else if conn, ok := msg.(connGenerated); ok {
		m.conn = conn
		m.msg = "Opened connection with Beepy…"
		return m, genClient(m.conn)
	} else if client, ok := msg.(clientGenerated); ok {
		m.client = client
		m.msg = "Created SFTP client…"
		return m, makeRemoteDirs(m.client, m.dirs)
	} else if _, ok := msg.(remoteDirsGenerated); ok {
		m.msg = "Created gomuks directories on Beepy…"
		return m, exfiltrate(m.client, m.patterns, m.dirs, m.root)
	} else if _, ok := msg.(exfilSuccess); ok {
		m.msg = "Copied gomuks instance to Beepy…"
		return m, downloadLatestGomuksBinary
	} else if archive, ok := msg.(gomuksFetched); ok {
		m.msg = "Downloaded latest gomuks binary…"
		return m, decompressGomuksDownload(string(archive))
	} else if binary, ok := msg.(gomuksDecompressed); ok {
		m.msg = "Unzipped gomuks binary, transfering to Beepy…"
		return m, transferGomuks(string(binary), m.client)
	} else if _, ok := msg.(gomuksTransfered); ok {
		m.status = success
	} else if err, ok := msg.(transferErr); ok {
		m.err = err
		m.status = failure
	} else if key, ok := msg.(tea.KeyMsg); ok {
		if key.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}

		switch m.status {
		case failure:
			m.err = nil
			m.msg = defaultMsg
			m.status = loading
			return m, utils.PrevPhase
		case success:
			return m, utils.NextPhase
		}
	} else if _, ok := msg.(spinner.TickMsg); ok {
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
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
			utils.Title().Render("Gomuks binary transfer failed"),
			utils.Error(m.err.Error()),
		)
	case success:
		return fmt.Sprintf(
			"%s\n"+
				"Well done, apprentice.\n"+
				"Press any key to continue.",
			utils.Title().Render("Gomuks binary transfer successful"),
		)
	}
	return ""
}
