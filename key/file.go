package key

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbletea"
)

type keyFileMsg string

func getTextFileInWd() tea.Msg {
	file := ""

	cwd, err := os.Getwd()
	if err != nil {
		return nil
	}

	err = filepath.WalkDir(cwd, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if filepath.Ext(d.Name()) == ".txt" {
			file = path
			return fs.SkipAll
		}

		return nil
	})
	if err != nil || file == "" {
		return nil
	}

	return keyFileMsg(file)
}
