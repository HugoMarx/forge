package projects

import (
	"fmt"
	"os"
	"time"

	"hugom/forge/config"

)

const (
	RootDir string = "/home/hugo/Projects"
)

type Project struct {
	Name     string
	Modified string
	DirSize  string
}

func (p Project) ToRow() []string {
	return []string{p.Name, p.Modified, p.DirSize}
}

func DiscoverProjects()(projects []Project,err error) {
	config, err := config.GetConfig()
	if err != nil {
		 return nil, fmt.Errorf("%w", err)
	}
	var discoveredProjects []Project
	entries, err := os.ReadDir(os.ExpandEnv(config.RootDir))
	if err != nil {
		fmt.Println(err)
	}

	for _, entry := range entries {
		entryInfo, err := entry.Info()
		if err != nil {
			fmt.Println(err)
		}

		// TODO Optimiser l'execution
		// dirSize, err := dirSize(RootDir + "/" + entry.Name())
		// if err != nil {
		// 	fmt.Println(err)
		// }

		discoveredProjects = append(discoveredProjects, Project{
			entry.Name(),
			entryInfo.ModTime().Format(time.DateOnly),
			"n/a", // dirSize,
		})
	}
	return discoveredProjects, nil
}

// Temporairement inactive.
// func dirSize(dirPath string) (string, error) {
// 	const (
// 		KB = 1024
// 		MB = KB * 1024
// 		GB = MB * 1024
// 	)

// 	totalSize := 0
// 	err := filepath.Walk(dirPath, func(path string, info fs.FileInfo, err error) error {
// 		totalSize += int(info.Size())
// 		return nil
// 	})
// 	if err != nil {
// 		return fmt.Sprint(0), err
// 	}
// 	return fmt.Sprintf("%.2f GB", float64(totalSize)/float64(GB)), nil
// }
