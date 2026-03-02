package main

import (
	"fmt"
	"os"
)

func main() {
	files, err := os.ReadDir("/home/hugom/Projects")
	if err != nil {
		fmt.Println(err)
	}
	for _, file := range files {
		if file.IsDir() {
			fmt.Println(file.Name())
		}
	}
}
