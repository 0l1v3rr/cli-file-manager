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

	if fileStat.IsDir() {
		return fmt.Sprintf("[Name:](fg:green) %v\n[Path:](fg:green) %v\n[Size:](fg:green) %v\n[Permission:](fg:green) %v\n[Directory:](fg:green) %v\n[Last Modified:](fg:green) %v", fileStat.Name(), path, "Press space to calculate", fileStat.Mode(), isDir, fileStat.ModTime().Format(time.RFC1123))
	} else {
		return fmt.Sprintf("[Name:](fg:green) %v\n[Path:](fg:green) %v\n[Size:](fg:green) %v\n[Permission:](fg:green) %v\n[Directory:](fg:green) %v\n[Last Modified:](fg:green) %v", fileStat.Name(), path, formatByte(int(fileStat.Size())), fileStat.Mode(), isDir, fileStat.ModTime().Format(time.RFC1123))
	}
}

func GetFileInformationsWithSize(p string) string {
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

	return fmt.Sprintf("[Name:](fg:green) %v\n[Path:](fg:green) %v\n[Size:](fg:green) %v\n[Permission:](fg:green) %v\n[Directory:](fg:green) %v\n[Last Modified:](fg:green) %v", fileStat.Name(), path, formatByte(size), fileStat.Mode(), isDir, fileStat.ModTime().Format(time.RFC1123))
}

func formatByte(size int) string {
	if size <= 1000 {
		return fmt.Sprintf("%v byte", size)
	}

	if size/1000 <= 1000 {
		return fmt.Sprintf("%.1f kilobyte", float32(float32(size)/1000))
	}

	if size/1000/1000 <= 1000 {
		return fmt.Sprintf("%.1f megabyte", float32(float32(size)/1000/1000))
	}

	if size/1000/1000/1000 <= 1000 {
		return fmt.Sprintf("%.1f gigabyte", float32(float32(size)/1000/1000/1000))
	}

	return ""
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
