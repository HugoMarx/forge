package docker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"hugom/forge/forgemsg"
	"hugom/forge/projects"

	tea "charm.land/bubbletea/v2"
)

func DockerComposeInspect(projectName string, format string) tea.Cmd {
	return func() tea.Msg {
		if format == "" {
			format = "{{json .}}"
		}

		hasCompose, projectDir := HasDockerComposeFile(projectName)
		if !hasCompose {
			return NoDockerFileMsg{fmt.Sprintln("Aucune config Docker détectée dans", projectDir, "!")}
		}

		cmd := exec.Command("docker", "compose", "ps", "-a", fmt.Sprint("--format=", format))
		cmd.Dir = projectDir
		debugDir := fmt.Sprint("Running in dir:", cmd.Dir)
		debugCommand := fmt.Sprint("Command:", cmd.Args)
		output, err := cmd.CombinedOutput()
		if err != nil {
			return forgemsg.CmdErrorMsg{Error: err, Debug: []string{debugDir, debugCommand}}
		}
		lines := bytes.Split(output, []byte("\n"))
		var containers []Container
		for _, line := range lines {
			line = bytes.TrimSpace(line)
			if len(line) == 0 {
				continue
			}
			var c Container
			if err := json.Unmarshal(line, &c); err != nil {
				return forgemsg.CmdErrorMsg{Error: err, Debug: nil}
			}
			containers = append(containers, c)
		}

		return ContainerInspectMsg{Project: projectName, Containers: containers}
	}
}

func DockerComposeUp(projectName string, launch bool) tea.Cmd {
	DockerComposeDown(projectName) // On démonte systématiquement le container avant de le lancer.

	return func() tea.Msg {
		cmd := exec.Command("docker", "compose", "up", "-d")
		cmd.Dir = filepath.Join(projects.RootDir, projectName)

		output, err := cmd.CombinedOutput()
		if err != nil {
			return forgemsg.CmdErrorMsg{Error: err, Debug: []string{}}
		}

		return DockerStateMsg{Project: projectName, Error: nil, IsRunning: true, Output: output, Options: map[string]any{"launch": launch}}
	}
}

func DockerComposeDown(projectName string, options ...[]string) tea.Cmd {
	return func() tea.Msg {
		cmd := exec.Command("docker", "compose", "down")
		cmd.Dir = filepath.Join(projects.RootDir, projectName)
		output, err := cmd.CombinedOutput()
		if err != nil {
			return forgemsg.CmdErrorMsg{Error: err, Debug: []string{}}
		}
		return DockerStateMsg{Project: projectName, Error: nil, IsRunning: false, Output: output}
	}
}

func HasDockerComposeFile(projectName string) (bool, string) {
	projectDir := filepath.Join(projects.RootDir, projectName)
	_, err := os.Stat(fmt.Sprintf("%s/docker-compose.yml", projectDir))
	return !os.IsNotExist(err), projectDir
}
