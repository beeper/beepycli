package utils

import (
	"github.com/charmbracelet/bubbletea"
)

type NextPhaseMsg struct{}
type PrevPhaseMsg struct{}

func NextPhase() tea.Msg {
	return NextPhaseMsg{}
}

func PrevPhase() tea.Msg {
	return PrevPhaseMsg{}
}
