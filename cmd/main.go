package main

import (
	"fmt"
	"os"

	"charm.land/bubbles/v2/table"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"hugom/forge/internals"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type model struct {
	table  table.Model
	output viewport.Model
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Une erreur est survenue: %v", err)
		os.Exit(1)
	}
}

func initialModel() model {
	columns := []table.Column{
		{Title: "Projet", Width: 60},
		{Title: "Modified", Width: 10},
		{Title: "Size", Width: 20},
	}

	rows := []table.Row{}

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")). //#ff4d43 Forge Red
		Bold(true)
	discoveredProjects := internals.DiscoverProjects()
	for _, discoveredProject := range discoveredProjects {
		rows = append(rows, table.Row{discoveredProject.Name, discoveredProject.Modified, discoveredProject.DirSize})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(25),
	)

	v := viewport.New(
		viewport.WithHeight(20),
	)

	startContent, _ := os.ReadFile("assets/forge.txt")
	v.SetContent(" Welcome to" + fmt.Sprintf("\x1b[31m%s\x1b[0m\n", string(startContent)))

	t.SetStyles(s)

	return model{t, v}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:

		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "enter":
			project := m.table.SelectedRow()[0]
			return m, internals.RunCommand(project, "ghostty", "-e", "hx", internals.RootDir+"/"+project)
		case "j":
			m.output.HalfPageDown()
			return m, nil
		case "k":
			m.output.HalfPageUp()
			return m, nil
		}

	case tea.WindowSizeMsg:
		totalWidth := msg.Width - 7 // Marge pour la scrollbar

		columns := []table.Column{
			{Title: "Projet", Width: int(float64(totalWidth) * 0.8)},
			{Title: "Modified", Width: int(float64(totalWidth) * 0.1)},
			{Title: "Size", Width: int(float64(totalWidth) * 0.1)},
		}

		m.table.SetColumns(columns)
		m.table.SetWidth(totalWidth)
		m.table.SetHeight(msg.Height / 2)
		m.output.SetWidth(msg.Width - 2)
		m.output.SetHeight(msg.Height / 3)

	case internals.CommandSuccessMsg:
		content := m.output.GetContent()
		content += fmt.Sprintf("Projet %s ouvert dans votre IDE préféré !\n", msg.ProjectName)
		m.output.SetContent(content)
		return m, nil

	}
	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) View() tea.View {
	rows := m.table.Rows()
	renderedRows := make([]table.Row, len(rows))
	for i, row := range rows {
		firstCell := row[0]
		if i == m.table.Cursor() {
			firstCell = "-> " + firstCell
		}
		renderedRows[i] = table.Row{firstCell, row[1], row[2]}
	}

	m.table.SetRows(renderedRows)

	tableRender := baseStyle.Render(m.table.View())
	outputRender := baseStyle.Render(m.output.View())

	return tea.NewView(lipgloss.JoinVertical(lipgloss.Top, tableRender, outputRender))
}
