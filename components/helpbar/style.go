package helpbar

import "charm.land/lipgloss/v2"

func GetStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color("201"))
}
