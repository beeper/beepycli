package transfer

import (
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbletea"
	"gopkg.in/yaml.v3"
)

type gomuksConfig struct {
	UserID      string `yaml:"mxid"`
	DeviceID    string `yaml:"device_id"`
	AccessToken string `yaml:"access_token"`
	HS          string `yaml:"homeserver"`

	RoomCacheSize int   `yaml:"room_cache_size"`
	RoomCacheAge  int64 `yaml:"room_cache_age"`

	NotifySound        bool `yaml:"notify_sound"`
	SendToVerifiedOnly bool `yaml:"send_to_verified_only"`

	Backspace1RemovesWord bool `yaml:"backspace1_removes_word"`
	Backspace2RemovesWord bool `yaml:"backspace2_removes_word"`

	AlwaysClearScreen bool `yaml:"always_clear_screen"`

	DataDir      string `yaml:"data_dir"`
	CacheDir     string `yaml:"cache_dir"`
	HistoryPath  string `yaml:"history_path"`
	RoomListPath string `yaml:"room_list_path"`
	MediaDir     string `yaml:"media_dir"`
	DownloadDir  string `yaml:"download_dir"`
	StateDir     string `yaml:"state_dir"`
}

type gomuksConfigUpdatedMsg struct{}

func updateGomuksConfigForNewLocation(username, root string, dirs *collection) tea.Cmd {
	return func() tea.Msg {
		path := filepath.Join(root, "config", "config.yaml")
		file, err := os.Open(path)
		if err != nil {
			return transferErr(err)
		}

		fields := gomuksConfig{}
		decoder := yaml.NewDecoder(file)
		err = decoder.Decode(&fields)
		if err != nil {
			return transferErr(err)
		}
		file.Close()

		fields.DataDir = dirs.data
		fields.CacheDir = dirs.cache
		fields.HistoryPath = filepath.Join(dirs.cache, "history.db")
		fields.RoomListPath = filepath.Join(dirs.cache, "rooms.gob.gz")
		fields.MediaDir = dirs.media
		fields.DownloadDir = filepath.Join("/home", username, "Downloads")
		fields.StateDir = dirs.state

		file, err = os.Create(path)
		if err != nil {
			return transferErr(err)
		}

		encoder := yaml.NewEncoder(file)
		err = encoder.Encode(fields)
		if err != nil {
			return transferErr(err)
		}
		encoder.Close()

		return gomuksConfigUpdatedMsg{}
	}
}
