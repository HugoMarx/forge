package helper

import tea "charm.land/bubbletea/v2"

type Layout struct {
	TermWidth                int
	TermHeigth               int
	LeftPanelWidth           int
	LeftPanelTopWinHeight    int
	LeftPanelBottomWinHeight int
	RightPanelWidth          int
	RightPanelWinHeight      int
}

func (layout *Layout) ComputeLayoutDimensions(msg tea.WindowSizeMsg) {
	layout.TermWidth = int(msg.Width)
	layout.TermHeigth = int(msg.Height)

	layout.LeftPanelWidth = int((msg.Width / 2) + 5)
	layout.RightPanelWidth = int((int(msg.Width) - layout.LeftPanelWidth) - 5)

	layout.LeftPanelTopWinHeight = int(msg.Height / 2)
	layout.LeftPanelBottomWinHeight = int(int(msg.Height)-layout.LeftPanelTopWinHeight) - (int(float64(msg.Height) * 0.14))
	layout.RightPanelWinHeight = int(msg.Height) - (int(float64(msg.Height) * 0.2))
}
