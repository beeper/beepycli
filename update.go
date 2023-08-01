package main

import (
	"github.com/charmbracelet/bubbletea"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyEnter, tea.KeyTab:
			switch m.phase {
			case welcome:
				if !m.focusOnButton {
					m.focusOnButton = true
				} else {
					m.phase = ssh
					m.focusOnButton = false
				}
			case ssh:
				if m.ssh.host.Focused() {
					m.focusOnButton = true
					m.ssh.host.Blur()
				} else if m.ssh.username.Focused() {
					m.ssh.username.Blur()
					m.ssh.username.CursorStart()
					m.ssh.host.Focus()
				} else if m.focusOnButton {
					m.phase = matrix
					m.focusOnButton = false
				} else {
					m.ssh.username.CursorEnd()
					m.ssh.username.Focus()
				}
			case matrix:
				if m.matrix.homeserver.Focused() {
					m.focusOnButton = true
					m.matrix.homeserver.Blur()
				} else if m.matrix.password.Focused() {
					m.matrix.password.Blur()
					m.matrix.homeserver.Focus()
				} else if m.matrix.username.Focused() {
					m.matrix.username.Blur()
					m.matrix.password.Focus()
				} else if m.focusOnButton {
					m.phase = next
					m.focusOnButton = false
				} else {
					m.matrix.username.Focus()
				}
			}
		case tea.KeyShiftTab:
			switch m.phase {
			case welcome:
				m.focusOnButton = false
			case ssh:
				if m.focusOnButton {
					m.focusOnButton = false
					m.ssh.host.Focus()
				} else if m.ssh.host.Focused() {
					m.ssh.username.CursorEnd()
					m.ssh.host.Blur()
					m.ssh.username.Focus()
				} else if m.ssh.username.Focused() {
					m.ssh.username.Blur()
					m.ssh.username.CursorStart()
				} else {
					m.phase = welcome
					m.focusOnButton = true
				}
			case matrix:
				if m.focusOnButton {
					m.focusOnButton = false
					m.matrix.homeserver.Focus()
				} else if m.matrix.homeserver.Focused() {
					m.matrix.homeserver.Blur()
					m.matrix.password.Focus()
				} else if m.matrix.password.Focused() {
					m.matrix.password.Blur()
					m.matrix.username.Focus()
				} else if m.matrix.username.Focused() {
					m.matrix.username.Blur()
				} else {
					m.phase = ssh
					m.focusOnButton = true
				}
			}
		default:
			var cmd tea.Cmd
			if m.phase == ssh {
				if m.ssh.username.Focused() {
					m.ssh.username, cmd = m.ssh.username.Update(msg)
				} else if m.ssh.host.Focused() {
					m.ssh.host, cmd = m.ssh.host.Update(msg)
				}
			} else if m.phase == matrix {
				if m.matrix.homeserver.Focused() {
					m.matrix.homeserver, cmd = m.matrix.homeserver.Update(msg)
				} else if m.matrix.password.Focused() {
					m.matrix.password, cmd = m.matrix.password.Update(msg)
				} else if m.matrix.username.Focused() {
					m.matrix.username, cmd = m.matrix.username.Update(msg)
				}
			}
			return m, cmd
		}
	}
	return m, nil
}
