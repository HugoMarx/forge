package main

import (
	"fmt"
	"os"
	"strings"

	c "hugom/forge/components"
	"hugom/forge/components/forgetable"
	"hugom/forge/components/helpbar"
	"hugom/forge/docker"
	"hugom/forge/forgemsg"
	"hugom/forge/helper"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
)

type rootModel struct {
	projectsTable forgetable.ForgeTable
	commandOutput viewport.Model
	topWindow     viewport.Model
	centerWindow  viewport.Model
	bottomWindow  viewport.Model
	helpBar       help.Model
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Une erreur est survenue: %v", err)
		os.Exit(1)
	}
}

var stringBuilder strings.Builder

func initialModel() rootModel {
	commandOutput := viewport.New()
	forgeTitle, _ := os.ReadFile("assets/forge.txt")
	// TODO Utiliser les color helpers lipgloss
	startupMessage := " Welcome to" + fmt.Sprintf("\x1b[31m%s\x1b[0m\n", forgeTitle)
	commandOutput.SetContent(startupMessage)

	return rootModel{
		projectsTable: forgetable.ProjectsTable,
		commandOutput: commandOutput,
		topWindow:     viewport.New(),
		centerWindow:  viewport.New(),
		bottomWindow:  viewport.New(),
		helpBar:       help.New(),
	}
}

func (m rootModel) Init() tea.Cmd {
	return nil
}

func (m rootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	selectedProject := m.projectsTable.Table.SelectedRow()[0]
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyPressMsg:

		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "enter":
			return m, helper.LaunchWorkspace(selectedProject)
		case "u":
			return m, docker.DockerComposeUp(selectedProject, false)
		case "U":
			return m, docker.DockerComposeUp(selectedProject, true)
		case "d":
			stringBuilder.Reset()
			fmt.Fprintf(&stringBuilder, "%s shutting down ! 🐋\n", selectedProject)
			m.topWindow.SetContent(stringBuilder.String())
			return m, docker.DockerComposeDown(selectedProject)
		case "s":
			return m, docker.DockerComposeInspect(selectedProject, "")
		case "up", "down":
			m.projectsTable.Table, cmd = m.projectsTable.Table.Update(msg)
			selected := m.projectsTable.Table.SelectedRow()[0]
			return m, tea.Batch(docker.DockerComposeInspect(selected, ""), cmd)
		case "k":
			m.commandOutput.HalfPageUp()
			return m, nil
		case "?":
			m.commandOutput.SetContent(m.helpBar.FullHelpView(helpbar.Keys.FullHelp()))
			return m, nil
		}

	case tea.WindowSizeMsg:
		leftPanelWidth := (msg.Width / 2) + 5
		rightPanelWidth := (msg.Width - leftPanelWidth) - 5

		leftPanelTopWinHeight := msg.Height / 2
		leftPanelBottomWinHeight := float64(msg.Height-leftPanelTopWinHeight) - (float64(msg.Height) * float64(0.1))
		rightPanelWinHeight := float64(msg.Height) - (float64(msg.Height) * float64(0.1))

		m.projectsTable.ResizeColumns(msg.Width)
		m.projectsTable.Table.SetWidth(leftPanelWidth)
		m.projectsTable.Table.SetHeight(leftPanelTopWinHeight)

		m.commandOutput.SetWidth(leftPanelWidth)
		m.commandOutput.SetHeight(int(leftPanelBottomWinHeight))

		for key, rigthPanelWindow := range []*viewport.Model{&m.topWindow, &m.centerWindow, &m.bottomWindow} {
			rigthPanelWindow.SetHeight(int((rightPanelWinHeight * 0.25)))
			if key == 0 {
				rigthPanelWindow.SetHeight(int(rightPanelWinHeight * 0.5))
			}
			rigthPanelWindow.SetWidth(rightPanelWidth)
		}

	case forgemsg.CmdSuccessMsg:
		outputContent := msg.Output
		m.commandOutput.SetContent(outputContent)
		m.commandOutput.GotoBottom()
		return m, nil
	case docker.NoDockerFileMsg:
		stringBuilder.Reset()
		fmt.Fprint(&stringBuilder, msg.Message)
		m.topWindow.SetContent(ansi.Wordwrap(stringBuilder.String(), m.topWindow.Width(), ""))
		return m, nil
	case docker.ContainerInspectMsg:
		stringBuilder.Reset()

		if len(msg.Containers) == 0 {
			fmt.Fprint(&stringBuilder, "Aucun conteneur monté :\nTapez u pour initialiser l'environnement Docker.")
		} else {
			fmt.Fprintf(&stringBuilder, "Docker data pour %s\n", msg.Project)
			for _, container := range msg.Containers {
				fmt.Fprintf(&stringBuilder, "\t%s\n\t%s\n\t%s\n\t%s\n\t%s\n\n", container.Name, container.Status, container.Service, container.Ports, container.State)
			}

		}

		m.topWindow.SetContent(stringBuilder.String())
		return m, nil
	case docker.ContainerStateMsg:
		var cmd tea.Cmd
		if !msg.IsRunning {
			fmt.Fprint(&stringBuilder, string(msg.Output))
		} else {
			stringBuilder.Reset()
			if msg.IsRunning {
				fmt.Fprint(&stringBuilder, string(msg.Output))
				fmt.Fprintf(&stringBuilder, "%s is now running in Docker ! 🐋\n", msg.Project)

				if value, ok := msg.Options["launch"]; ok {
					if withLaunch, ok := value.(bool); ok && withLaunch {
						cmd = helper.LaunchWorkspace(msg.Project)
					}
				}
			}
		}
		m.topWindow.SetContent(stringBuilder.String())
		return m, cmd
	case forgemsg.CmdErrorMsg:
		outputContent := m.commandOutput.GetContent()
		outputContent += fmt.Sprintln(msg.Error.Error())
		m.commandOutput.SetContent(outputContent + "\n" + strings.Join(msg.Debug, "\n") + "\n\n\n")
		m.commandOutput.GotoBottom()
		return m, nil
	}
	m.projectsTable.Table, cmd = m.projectsTable.Table.Update(msg)
	return m, cmd
}

func (m rootModel) View() tea.View {
	tableRender := m.projectsTable.Render()
	outputRender := c.BaseStyle.Render(m.commandOutput.View())
	topMonitoringRender := c.BaseStyle.Render(m.topWindow.View())
	centerMonitoringRender := c.BaseStyle.Render(m.centerWindow.View())
	bottomMonitoringRender := c.BaseStyle.Render(m.bottomWindow.View())
	helpRender := helpbar.GetStyle().Render(m.helpBar.ShortHelpView(helpbar.Keys.ShortHelp()))

	leftComposedViewport := lipgloss.JoinVertical(lipgloss.Left, tableRender, outputRender)
	rightComposedViewport := lipgloss.JoinVertical(lipgloss.Left, topMonitoringRender, centerMonitoringRender, bottomMonitoringRender)
	return tea.NewView(
		lipgloss.JoinHorizontal(lipgloss.Top, leftComposedViewport, rightComposedViewport) + "\n" + helpRender)
}
