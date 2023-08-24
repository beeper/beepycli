package matrix

import (
	"maunium.net/go/gomuks/beeper"

	"github.com/charmbracelet/bubbletea"
)

type apiError error
type loginStarted string
type emailSuccess struct{}

func initAuth() tea.Cmd {
	return func() tea.Msg {
		resp, err := beeper.StartLogin()
		if err != nil {
			return apiError(err)
		}
		return loginStarted(resp.RequestID)
	}
}

func sendEmail(session, email string) tea.Cmd {
	return func() tea.Msg {
		err := beeper.SendLoginEmail(session, email)
		if err != nil {
			return apiError(err)
		}
		return emailSuccess{}
	}
}
