package main

import (
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/id"

	"github.com/charmbracelet/bubbletea"
)

type usernameParseMsg string
type usernameErrMsg string
type homeserverParseMsg string
type homeserverErrMsg string

func getHomeserverString(username string) tea.Cmd {
	return func() tea.Msg {
		_, homeserver, err := id.UserID(username).Parse()
		if err != nil {
			return usernameErrMsg(err.Error())
		}
		return usernameParseMsg(homeserver)
	}
}

func resolveWellKnown(homeserver string) tea.Cmd {
	return func() tea.Msg {
		resp, err := mautrix.DiscoverClientAPI(homeserver)
		if err != nil {
			return homeserverErrMsg(err.Error())
		} else if resp != nil {
			return homeserverParseMsg(resp.Homeserver.BaseURL)
		}
		return homeserverErrMsg("Error: parsing homeserver failed")
	}
}
