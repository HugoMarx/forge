package main

import (
	"fmt"
	"os"
	"strings"

	"hugom/forge/components/forgetable"
	"hugom/forge/components/helpbar"
	"hugom/forge/docker"
	"hugom/forge/forgemsg"
	"hugom/forge/helper"
	"hugom/forge/projects"

	c "hugom/forge/components"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
)

type rootModel struct {
	projectsTable *forgetable.ForgeTable
	commandOutput viewport.Model
	dockerTable   *forgetable.ForgeTable
	topWindow     viewport.Model
	centerWindow  viewport.Model
	bottomWindow  viewport.Model
	spinner       spinner.Model
	helpBar       help.Model
	layout        helper.Layout
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

	discoveredProjects := projects.DiscoverProjects()
	forgetable.MainTable.BuildTable(forgetable.ToRowable(discoveredProjects), helper.Layout{})
	return rootModel{
		projectsTable: forgetable.MainTable,
		commandOutput: commandOutput,
		dockerTable:   forgetable.DockerTable,
		topWindow:     viewport.New(),
		centerWindow:  viewport.New(),
		bottomWindow:  viewport.New(),
		spinner: spinner.New(
			spinner.WithSpinner(spinner.Pulse),
			spinner.WithStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("222")))),
		helpBar: help.New(),
	}
}

func (m rootModel) Init() tea.Cmd {
	return m.spinner.Tick
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
			hasCompose, _ := docker.HasDockerComposeFile(selectedProject)
			if !hasCompose {
				return m, nil
			}
			m.dockerTable.IsLoading = true
			return m, docker.DockerComposeUp(selectedProject, false)
		case "U":
			hasCompose, _ := docker.HasDockerComposeFile(selectedProject)
			if !hasCompose {
				return m, nil
			}

			m.dockerTable.IsLoading = true
			return m, docker.DockerComposeUp(selectedProject, true)
		case "d":
			hasCompose, _ := docker.HasDockerComposeFile(selectedProject)
			if !hasCompose {
				return m, nil
			}
			m.dockerTable.IsLoading = true
			m.dockerTable.HasData = false
			stringBuilder.Reset()
			fmt.Fprintf(&stringBuilder, "%s shutting down ! 🐋\n", selectedProject)
			m.commandOutput.SetContent(stringBuilder.String())
			return m, docker.DockerComposeDown(selectedProject)
		case "s":
			hasCompose, _ := docker.HasDockerComposeFile(selectedProject)
			if !hasCompose {
				return m, nil
			}
			return m, docker.DockerComposeInspect(selectedProject, "")
		case "up", "down":
			m.projectsTable.Table, _ = m.projectsTable.Table.Update(msg)
			selected := m.projectsTable.Table.SelectedRow()[0]
			return m, docker.DockerComposeInspect(selected, "")
		case "k":
			m.commandOutput.HalfPageUp()
			return m, nil
		case "?":
			m.commandOutput.SetContent(m.helpBar.FullHelpView(helpbar.Keys.FullHelp()))
			return m, nil
		}

	case tea.WindowSizeMsg:
		m.layout.ComputeLayoutDimensions(msg)

		m.projectsTable.ResizeColumns(msg.Width)
		m.dockerTable.ResizeColumns(msg.Width)
		m.projectsTable.Table.SetWidth(m.layout.LeftPanelWidth)
		m.projectsTable.Table.SetHeight(m.layout.LeftPanelTopWinHeight)

		m.commandOutput.SetWidth(m.layout.LeftPanelWidth)
		m.commandOutput.SetHeight(m.layout.LeftPanelBottomWinHeight)

		for key, rigthPanelWindow := range []*viewport.Model{&m.topWindow, &m.centerWindow, &m.bottomWindow} {
			rigthPanelWindow.SetHeight(int(float64(m.layout.RightPanelWinHeight) * 0.2))
			if key == 0 {
				rigthPanelWindow.SetHeight(int(float64(m.layout.RightPanelWinHeight) * 0.5))
			}
			rigthPanelWindow.SetWidth(m.layout.RightPanelWidth)
		}

	case forgemsg.CmdSuccessMsg:
		outputContent := msg.Output
		m.commandOutput.SetContent(outputContent)
		m.commandOutput.GotoBottom()
		return m, cmd
	case docker.NoDockerFileMsg:
		stringBuilder.Reset()
		m.dockerTable.HasData = false
		fmt.Fprint(&stringBuilder, msg.Message)
		m.topWindow.SetContent(ansi.Wordwrap(stringBuilder.String(), m.topWindow.Width(), ""))
		return m, nil
	case docker.ContainerInspectMsg:
		stringBuilder.Reset()
		m.dockerTable.BuildTable(forgetable.ToRowable(msg.Containers), m.layout)
		if !m.dockerTable.HasData {
			fmt.Fprint(&stringBuilder, "Aucun conteneur monté :\nTapez u pour initialiser l'environnement Docker.")
			m.topWindow.SetContent(ansi.Wordwrap(stringBuilder.String(), m.topWindow.Width(), ""))
		}
		return m, nil
	case docker.DockerStateMsg:
		cmd = docker.DockerComposeInspect(msg.Project, "")
		if !msg.IsRunning {
			fmt.Fprint(&stringBuilder, string(msg.Output))
		} else {
			stringBuilder.Reset()
			fmt.Fprint(&stringBuilder, string(msg.Output))
			fmt.Fprintf(&stringBuilder, "%s is now running in Docker ! 🐋\n", msg.Project)

			if value, ok := msg.Options["launch"]; ok {
				if withLaunch, ok := value.(bool); ok && withLaunch {
					cmd = helper.LaunchWorkspace(msg.Project)
				}
			}
			m.dockerTable.BuildTable(forgetable.ToRowable(msg.Containers), m.layout)
		}
		m.commandOutput.SetContent(stringBuilder.String())
		m.dockerTable.IsLoading = false
		return m, cmd
	case forgemsg.CmdErrorMsg:
		outputContent := fmt.Sprintln(msg.Error.Error())
		m.commandOutput.SetContent(outputContent)
		m.commandOutput.GotoBottom()
		return m, nil
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		if m.dockerTable.IsLoading {
			m.topWindow.SetContent(" " + m.spinner.View() + " Loading...")
		}
		return m, cmd
	}
	m.projectsTable.Table, cmd = m.projectsTable.Table.Update(msg)
	return m, cmd
}

func (m rootModel) View() tea.View {
	outputRender := c.BaseStyle.Render(m.commandOutput.View())
	var topMonitoringRender string
	if m.dockerTable.HasData {
		s := forgetable.GetStyle()
		s.Selected = lipgloss.NewStyle() // style vide = pas de highlight
		m.dockerTable.Table.SetStyles(s)
		topMonitoringRender = m.dockerTable.Render()
	} else {
		topMonitoringRender = c.BaseStyle.Render(m.topWindow.View())
	}
	tableRender := m.projectsTable.Render()
	centerMonitoringRender := c.BaseStyle.Render(m.centerWindow.View())
	bottomMonitoringRender := c.BaseStyle.Render(m.bottomWindow.View())
	helpRender := helpbar.GetStyle().Render(m.helpBar.ShortHelpView(helpbar.Keys.ShortHelp()))

	leftComposedViewport := lipgloss.JoinVertical(lipgloss.Left, tableRender, outputRender)
	rightComposedViewport := lipgloss.JoinVertical(lipgloss.Left, topMonitoringRender, centerMonitoringRender, bottomMonitoringRender)
	return tea.NewView(
		lipgloss.JoinHorizontal(lipgloss.Top, leftComposedViewport, rightComposedViewport) + "\n" + helpRender)
}
