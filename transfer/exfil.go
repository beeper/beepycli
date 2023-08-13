package transfer

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"

	"github.com/charmbracelet/bubbletea"
)

type configGenerated *ssh.ClientConfig
type dirCollectionGenerated *collection
type patternCollectionGenerated *collection
type connGenerated *ssh.Client
type clientGenerated *sftp.Client
type remoteDirsGenerated struct{}
type exfilSuccess struct{}

func genDirCollection(username string) tea.Cmd {
	return func() tea.Msg {
		cacheDir := filepath.Join("/home", username, ".cache", "gomuks")
		return dirCollectionGenerated(&collection{
			conf:  filepath.Join("/home", username, ".config", "gomuks"),
			data:  filepath.Join("/home", username, ".local", "share", "gomuks"),
			cache: cacheDir,
			state: filepath.Join(cacheDir, "state"),
			media: filepath.Join(cacheDir, "media"),
		})
	}
}

func genPatternCollection(root string) tea.Cmd {
	return func() tea.Msg {
		return patternCollectionGenerated(&collection{
			conf:  filepath.Join(root, "config", "*"),
			data:  filepath.Join(root, "data", "gomuks", "*"),
			cache: filepath.Join(root, "cache", "*"),
			state: filepath.Join(root, "cache", "state", "*"),
			media: filepath.Join(root, "cache", "media", "*", "*"),
		})
	}
}

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

func genClient(conn *ssh.Client) tea.Cmd {
	return func() tea.Msg {
		client, err := sftp.NewClient(conn)
		if err != nil {
			return transferErr(err)
		}

		return clientGenerated(client)
	}
}

func makeRemoteDirs(client *sftp.Client, dirs *collection) tea.Cmd {
	return func() tea.Msg {
		err := client.MkdirAll(dirs.conf)
		if err != nil {
			return transferErr(
				fmt.Errorf("Failed to create config directory on Beepy with error: %s", err),
			)
		}

		err = client.MkdirAll(dirs.cache)
		err = client.MkdirAll(dirs.state)
		err = client.MkdirAll(dirs.media)
		if err != nil {
			return transferErr(
				fmt.Errorf("Failed to create cache directory on Beepy with error: %s", err),
			)
		}

		err = client.MkdirAll(dirs.data)
		if err != nil {
			return transferErr(
				fmt.Errorf("Failed to create data directory on Beepy with error: %s", err),
			)
		}

		return remoteDirsGenerated{}
	}
}

func exfiltrate(client *sftp.Client, patterns, dirs *collection, root string) tea.Cmd {
	return func() tea.Msg {
		err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				return nil
			}

			var remote *sftp.File

			local, err := os.Open(path)
			if err != nil {
				return err
			}
			defer local.Close()

			if isConfig, _ := filepath.Match(patterns.conf, path); isConfig {
				remote, err = client.Create(filepath.Join(dirs.conf, d.Name()))
			} else if isData, _ := filepath.Match(patterns.data, path); isData {
				remote, err = client.Create(filepath.Join(dirs.data, d.Name()))
			} else if isCache, _ := filepath.Match(patterns.cache, path); isCache {
				remote, err = client.Create(filepath.Join(dirs.cache, d.Name()))
			} else if isState, _ := filepath.Match(patterns.state, path); isState {
				remote, err = client.Create(filepath.Join(dirs.state, d.Name()))
			} else if isMedia, _ := filepath.Match(patterns.media, path); isMedia {
				dir := filepath.Base(filepath.Dir(path))
				remote, err = client.Create(filepath.Join(dirs.media, dir, d.Name()))
			} else {
				return fmt.Errorf("File found in Gomuks root that doesn't conform to known patterns: %s", path)
			}
			if err != nil {
				return err
			}
			defer remote.Close()

			_, err = io.Copy(remote, local)
			if err != nil {
				return err
			}

			return nil
		})
		if err != nil {
			return transferErr(err)
		}

		return exfilSuccess{}
	}
}
