package settings

import (
	"github.com/charmbracelet/lipgloss/v2"
)

var (
	purple      = lipgloss.Color("99")
	lightPurple = lipgloss.Color("98")
	yellow      = lipgloss.Color("220")
	gray        = lipgloss.Color("245")
	lightGray   = lipgloss.Color("241")

	Header = lipgloss.NewStyle().Foreground(purple).Bold(true)
	Accent = lipgloss.NewStyle().Foreground(lightPurple)
	Base   = lipgloss.NewStyle().Foreground(lightGray)
	Alert  = lipgloss.NewStyle().Foreground(yellow).Bold(true)
	Muted  = lipgloss.NewStyle().Foreground(gray)
)

func DisableColors() {
	Header = lipgloss.NewStyle()
	Accent = lipgloss.NewStyle()
	Base = lipgloss.NewStyle()
	Alert = lipgloss.NewStyle()
	Muted = lipgloss.NewStyle()
}
