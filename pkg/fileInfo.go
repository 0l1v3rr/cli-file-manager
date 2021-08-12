package pkg

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func GetFileInformations(p string) string {
	fileStat, err := os.Stat(p)
	if err != nil {
		log.Fatal(err)
	}

	isDir := "No"
	if fileStat.IsDir() {
		isDir = "Yes"
	}

	path := p
	if fileStat.Name() == "../" || fileStat.Name() == ".." {
		splitted := strings.Split(path, "/")
		if len(splitted) > 0 {
			splitted = splitted[:len(splitted)-1]
		}
		path = strings.Join(splitted, "/")
	}
	path = strings.ReplaceAll(path, "//", "/")

	return fmt.Sprintf("[Name:](fg:green) %v\n[Path:](fg:green) %v\n[Size:](fg:green) %v byte\n[Permission:](fg:green) %v\n[Directory:](fg:green) %v\n[Last Modified:](fg:green) %v", fileStat.Name(), path, fileStat.Size(), fileStat.Mode(), isDir, fileStat.ModTime())
}

func EmptyFileInfo() string {
	return "[Name:](fg:green) -\n[Path:](fg:green) -\n[Size:](fg:green) 0 byte\n[Permission:](fg:green) -\n[Directory:](fg:green) -\n[Last Modified:](fg:green) -"
}
