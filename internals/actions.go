package internals

import (
	"fmt"
	"os/exec"

	tea "charm.land/bubbletea/v2"
)

type CommandOutputMsg struct {
	OutputContent string
	ProjectName   string
	ProjectPath   string
}

func RunCommand(project string, command string, args ...string) tea.Cmd {
	return func() tea.Msg {
		cmd := exec.Command(command, args...)
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("error while executing specified command : %v", err)
		}
		path := RootDir + "/" + project
		return CommandOutputMsg{string(output) , project , path}
	}
}
