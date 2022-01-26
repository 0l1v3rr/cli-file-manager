package pkg

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
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

	size := 0
	if fileStat.Name() == "../" || fileStat.Name() == ".." {
		size = 0
	} else if fileStat.IsDir() {
		size = int(dirSize(p))
	} else {
		size = int(fileStat.Size())
	}

	return fmt.Sprintf("[Name:](fg:green) %v\n[Path:](fg:green) %v\n[Size:](fg:green) %v byte\n[Permission:](fg:green) %v\n[Directory:](fg:green) %v\n[Last Modified:](fg:green) %v", fileStat.Name(), path, size, fileStat.Mode(), isDir, fileStat.ModTime().Format(time.RFC1123))
}

func dirSize(path string) int64 {
	var size int64
	filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return size
}

func EmptyFileInfo() string {
	return "[Name:](fg:green) -\n[Path:](fg:green) -\n[Size:](fg:green) 0 byte\n[Permission:](fg:green) -\n[Directory:](fg:green) -\n[Last Modified:](fg:green) -"
}
