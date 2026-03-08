package main

import (
	"fmt"
	"os"
	"strings"

	"charm.land/bubbles/v2/table"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"hugom/forge/internals"
	"hugom/forge/internals/docker"
	"hugom/forge/internals/projects"
)

type rootModel struct {
	width  int
	height int

	projectsTable table.Model
	output        viewport.Model
	monitoring    viewport.Model
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Une erreur est survenue: %v", err)
		os.Exit(1)
	}
}

func initialModel() rootModel {
	columns := []table.Column{
		{Title: "Projet", Width: 60},
		{Title: "Modified", Width: 10},
		{Title: "Size", Width: 20},
	}

	rows := []table.Row{}

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
	discoveredProjects := projects.DiscoverProjects()
	for _, discoveredProject := range discoveredProjects {
		rows = append(rows, table.Row{discoveredProject.Name, discoveredProject.Modified, discoveredProject.DirSize})
	}

	tableModel := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
	)

	outputModel := viewport.New()
	monitoringModel := viewport.New()

	initialContent, _ := os.ReadFile("assets/forge.txt")
	outputModel.SetContent(" Welcome to" + fmt.Sprintf("\x1b[31m%s\x1b[0m\n", string(initialContent)))

	tableModel.SetStyles(tableStyle)

	return rootModel{0, 0, tableModel, outputModel, monitoringModel}
}

func (m rootModel) Init() tea.Cmd {
	return nil
}

func (m rootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	project := m.projectsTable.SelectedRow()[0]

	switch msg := msg.(type) {
	case tea.KeyPressMsg:

		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "enter":
			return m, internals.LaunchWorkspace(project)
		case "s":
			return m, internals.DockerComposeInspect(project, "")
		case "j":
			m.output.HalfPageDown()
			return m, nil
		case "k":
			m.output.HalfPageUp()
			return m, nil
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		leftWidth := m.width / 2
		rightWidth := m.width - leftWidth

		topHeight := m.height / 2
		bottomHeight := m.height - topHeight
		columns := []table.Column{
			{Title: "Projet", Width: int(float64(m.width/2) * 0.75)},
			{Title: "Modified", Width: int(float64(m.width/2) * 0.1)},
			{Title: "Size", Width: int(float64(m.width/2) * 0.1)},
		}
		m.projectsTable.SetColumns(columns)
		m.projectsTable.SetWidth(leftWidth)
		m.projectsTable.SetHeight(topHeight)
		m.output.SetWidth(leftWidth)
		m.output.SetHeight(bottomHeight - 5)
		m.monitoring.SetHeight(m.height - 3)
		m.monitoring.SetWidth(rightWidth - 10)

	case internals.CmdSuccessMsg:
		content := msg.Output
		m.output.SetContent(content)
		m.output.GotoBottom()
		return m, nil
	case docker.ContainerMsg:
		var builder strings.Builder
		builder.WriteString("Docker data pour " + msg.Project + "\n")
		for _, container := range msg.Containers {
			builder.WriteString("\t" + container.Name + "\n")
			builder.WriteString("\t" + container.Status + "\n")
			builder.WriteString("\t" + container.Service + "\n")
			builder.WriteString("\t" + container.Ports + "\n")
			builder.WriteString("\t" + container.State + "\n\n")
		}

		m.monitoring.SetContent(builder.String())
		return m, nil
	case internals.CmdErrorMsg:
		content := m.output.GetContent()
		content += fmt.Sprintln(msg.Error.Error())
		m.output.SetContent(content + "\n" + strings.Join(msg.Debug, "\n") + "\n\n\n")
		m.output.GotoBottom()
		return m, nil
	}
	var cmd tea.Cmd
	m.projectsTable, cmd = m.projectsTable.Update(msg)
	return m, cmd
}

func (m rootModel) View() tea.View {
	baseStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240"))
	tableRender := baseStyle.Render(m.projectsTable.View())
	outputRender := baseStyle.Render(m.output.View())
	monitoringRender := baseStyle.Render(m.monitoring.View())

	leftComposedViewport := lipgloss.JoinVertical(lipgloss.Left, tableRender, outputRender)
	return tea.NewView(
		lipgloss.JoinHorizontal(lipgloss.Center, leftComposedViewport, monitoringRender))
}
