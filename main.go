package main

import (
	"github.com/charmbracelet/bubbletea"
	gloss "github.com/charmbracelet/lipgloss"
)

type phase int

const (
	welcome phase = iota
	ssh
	matrix
	next

	passwordPlaceholder = "correct horse battery staple"
	domainPlaceholder   = "https://example.com"
)

var (
	magenta = gloss.Color("13")
	purple  = gloss.Color("5")
)

func main() {
	m := initModel()
	tea.NewProgram(m).Run()
}
