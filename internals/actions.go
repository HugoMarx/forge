package internals

import (
	"fmt"
	"os/exec"

	"hugom/forge/internals/projects"

	tea "charm.land/bubbletea/v2"
)

type CmdSuccessMsg struct {
	Output string
}

type CmdErrorMsg struct {
	Error error
}

func LaunchWorkspace(projectName string) tea.Cmd {
	return func() tea.Msg {
		cmd := exec.Command("wt.exe", "-M", "-p", "Ubuntu", "--title", projectName, "wsl", "bash", "-l", "-c", fmt.Sprintf("hx %s/%s ", projects.RootDir, projectName))
		err := cmd.Start()
		if err != nil {
			return CmdErrorMsg{Error: err}
		}
		return CmdSuccessMsg{Output: fmt.Sprintf("Projet %s ouvert dans votre IDE préféré !\n", projectName)}
	}
}
