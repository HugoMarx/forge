package internals

import (
	"os/exec"

	tea "charm.land/bubbletea/v2"
)

type CommandSuccessMsg struct {
	output any

}

type CmdErrorMsg struct {
	error error
}

func LauchWorkspace(projectPath string) tea.Msg {
	msg := RunCommand("ghossty", "-e", projectPath)
	msg.output = ""
	return msg
}

func RunCommand(command string, args ...string) tea.Cmd {
	return func() tea.Msg {
		cmd := exec.Command(command, args...)
		err := cmd.Start()
		if err != nil {
			return CmdErrorMsg{err}
		}

		return CommandSuccessMsg{}
	}
}
