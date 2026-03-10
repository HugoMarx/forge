package internals

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"

	"hugom/forge/internals/docker"
	"hugom/forge/internals/projects"

	tea "charm.land/bubbletea/v2"
)

type CmdSuccessMsg struct {
	Output string
}

type CmdErrorMsg struct {
	Error error
	Debug []string
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

func DockerComposeInspect(projectName string, format string) tea.Cmd {
	return func() tea.Msg {
		if format == "" {
			format = "{{json .}}"
		}
		cmd := exec.Command("docker", "compose", "ps", "-a", fmt.Sprint("--format=", format))
		cmd.Dir = filepath.Join(projects.RootDir, projectName)
		debugDir := fmt.Sprint("Running in dir:", cmd.Dir)
		debugCommand := fmt.Sprint("Command:", cmd.Args)
		output, err := cmd.CombinedOutput()
		if err != nil {
			return CmdErrorMsg{Error: err, Debug: []string{debugDir, debugCommand}}
		}
		lines := bytes.Split(output, []byte("\n"))
		var containers []docker.Container
		for _, line := range lines {
			line = bytes.TrimSpace(line)
			if len(line) == 0 {
				continue
			}
			var c docker.Container
			if err := json.Unmarshal(line, &c); err != nil {
				return CmdErrorMsg{err, nil}
			}
			containers = append(containers, c)
		}

		return docker.ContainerInspectMsg{Project: projectName, Containers: containers}
	}
}

func DockerComposeUp(projectName string, options ...[]string) tea.Cmd {
	return func() tea.Msg {
		cmd := exec.Command("docker", "compose", "up", "-d")
		cmd.Dir = filepath.Join(projects.RootDir, projectName)
		output, err := cmd.CombinedOutput()
		if err != nil {
			return docker.RunContainerMsg{Project: projectName, Error: err, IsRunning: false, Output: []byte{}}
		}
		return docker.RunContainerMsg{Project: projectName, Error: nil, IsRunning: false, Output: output}
	}
}
