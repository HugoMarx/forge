package forgetable

import (
	"charm.land/bubbles/v2/table"
	"charm.land/lipgloss/v2"
)

func GetStyle() table.Styles {
	tableStyle := table.DefaultStyles()
	tableStyle.Header = tableStyle.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)

	tableStyle.Selected = tableStyle.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("#ff4d43")). //#ff4d43 Forge Red
		Bold(true)
	return tableStyle
}
