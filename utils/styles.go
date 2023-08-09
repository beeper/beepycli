package utils

import gloss "github.com/charmbracelet/lipgloss"

func Title() gloss.Style {
	return gloss.NewStyle().Bold(true)
}

func Button(highlighted bool) gloss.Style {
	button := gloss.NewStyle().Background(Gray).Padding(0, 2)
	if highlighted {
		button = button.Background(Magenta)
	}
	return button
}
