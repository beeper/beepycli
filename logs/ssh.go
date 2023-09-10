package logs

import (
	"bytes"
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"

	"github.com/charmbracelet/bubbletea"
)

type transferErr error
type configGenerated *ssh.ClientConfig
type connGenerated *ssh.Client
type logFileFound string
type clientGenerated *sftp.Client
type logsCopiedSuccessfully string

func genConfig(username, password string) tea.Cmd {
	return func() tea.Msg {
		return configGenerated(&ssh.ClientConfig{
			User:            username,
			Auth:            []ssh.AuthMethod{ssh.Password(password)},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		})
	}
}

func genConn(host string, conf *ssh.ClientConfig) tea.Cmd {
	return func() tea.Msg {
		conn, err := ssh.Dial("tcp", host, conf)
		if err != nil {
			return transferErr(err)
		}

		return connGenerated(conn)
	}
}

func locateLogFileOnBeepy(conn *ssh.Client) tea.Cmd {
	return func() tea.Msg {
		session, err := conn.NewSession()
		if err != nil {
			return transferErr(err)
		}
		defer session.Close()

		var path bytes.Buffer
		session.Stdout = &path

		err = session.Run("gomuks --print-log-path")
		if err != nil {
			return transferErr(err)
		}

		return logFileFound(path.String())
	}
}

func genClient(conn *ssh.Client) tea.Cmd {
	return func() tea.Msg {
		client, err := sftp.NewClient(conn)
		if err != nil {
			return transferErr(err)
		}

		return clientGenerated(client)
	}
}

func copyLogs(client *sftp.Client, file string) tea.Cmd {
	return func() tea.Msg {
		remote, err := client.Open(file)
		if err != nil {
			return transferErr(err)
		}
		defer remote.Close()

		path, err := os.UserHomeDir()
		if err == nil {
			path = filepath.Join(path, "Downloads")
		}
		path = filepath.Join(path, "beepy.log")

		local, err := os.Create(path)
		if err != nil {
			return transferErr(err)
		}
		defer local.Close()

		_, err = io.Copy(local, remote)
		if err != nil {
			return transferErr(err)
		}

		return logsCopiedSuccessfully(path)
	}
}
