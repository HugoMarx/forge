package helper

import (
	"fmt"
	"hugom/forge/config"
	"hugom/forge/forgemsg"
	"os/exec"

	tea "charm.land/bubbletea/v2"
)

func LaunchWorkspace(projectName string) tea.Cmd {
	return func() tea.Msg {
		cmd := exec.Command("wt.exe", "-M", "-p", "Ubuntu", "--title", projectName, "wsl", "bash", "-l", "-c", fmt.Sprintf("hx %s/%s ", config.Config.RootDir, projectName))
		err := cmd.Start()
		if err != nil {
			return forgemsg.CmdErrorMsg{Error: err}
		}
		return forgemsg.CmdSuccessMsg{Output: fmt.Sprintf("Projet %s ouvert dans votre IDE préféré !\n", projectName)}
	}
}

