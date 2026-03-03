package internals

import (
	"fmt"
	"os/exec"

	"hugom/forge/internals/projects"

	tea "charm.land/bubbletea/v2"
)

type CommandSuccessMsg struct {
	ProjectName string
}

type CmdErrorMsg struct {
	Error error
}

func LaunchWorkspace(projectName string) tea.Cmd {
	return RunCommand(projectName, "wt.exe", "-M", "-p", "Ubuntu", "--title", projectName, "wsl", "bash", "-l", "-c", fmt.Sprintf("hx %s/%s ", projects.RootDir, projectName))
}

func RunCommand(projectName string, command string, args ...string) tea.Cmd {
	return func() tea.Msg {
		cmd := exec.Command(command, args...)
		err := cmd.Start()
		if err != nil {
			return CmdErrorMsg{err}
		}
		return CommandSuccessMsg{
			ProjectName: projectName,
		}
	}
}
