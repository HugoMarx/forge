package docker

type Container struct {
	Name    string `json:"Name"`
	Service string `json:"Service"`
	State   string `json:"State"`
	Status  string `json:"Status"`
	Ports   string `json:"Ports"`
}

func (c Container) ToRow() []string {
	return []string{
		c.Name,
		c.Service,
		c.State,
		c.Status,
		c.Ports,
	}
}

type ContainerInspectMsg struct {
	Project    string
	Error      error
	Containers []Container
	Output     []byte
}

type ContainerStateMsg struct {
	Project   string
	Error     error
	IsRunning bool
	Output    []byte
	Options   map[string]any
}

type NoDockerFileMsg struct {
	Message string
}
