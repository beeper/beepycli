package transfer

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/pkg/sftp"

	"github.com/charmbracelet/bubbletea"
)

const url = "https://github.com/beeper/gomuks/releases/download/nightly/gomuks-linux-armv7"

type gomuksFetched string
type gomuksTransfered struct{}

func downloadLatestGomuksBinary() tea.Msg {
	destination := filepath.Join(os.TempDir(), "gomuks")

	out, err := os.Create(destination)
	if err != nil {
		return transferErr(err)
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return transferErr(err)
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return transferErr(err)
	}

	return gomuksFetched(destination)
}

func transferGomuks(binary string, client *sftp.Client) tea.Cmd {
	return func() tea.Msg {
		local, err := os.Open(binary)
		if err != nil {
			return transferErr(err)
		}
		defer local.Close()

		remote, err := client.Create("/usr/local/bin/gomuks")
		if err != nil {
			return transferErr(err)
		}
		defer remote.Close()

		err = remote.Chmod(0744)
		if err != nil {
			return transferErr(err)
		}

		_, err = io.Copy(remote, local)
		if err != nil {
			return transferErr(err)
		}

		return gomuksTransfered{}
	}
}
