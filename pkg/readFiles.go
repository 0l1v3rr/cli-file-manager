package pkg

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

func ReadFiles(p string, showHidden bool) []string {
	files, err := ioutil.ReadDir(p)
	if err != nil {
		log.Fatal(err)
	}

	res := []string{}
	counter := 0

	if showHidden {
		if p != "/" {
			counter++
			res = append(res, fmt.Sprintf("[[%v] %v/](fg:green)", 1, ".."))
		}
		for _, file := range files {
			if file.IsDir() {
				counter++
				res = append(res, fmt.Sprintf("[[%v] %v/](fg:green)", counter, file.Name()))
			}
		}
		for _, file := range files {
			if !file.IsDir() {
				counter++
				res = append(res, fmt.Sprintf("[%v] %v", counter, file.Name()))
			}
		}
	} else {
		if p != "/" {
			counter++
			res = append(res, fmt.Sprintf("[[%v] %v/](fg:green)", 1, ".."))
		}
		for _, file := range files {
			if strings.HasPrefix(file.Name(), ".") {
				continue
			}
			if file.IsDir() {
				counter++
				res = append(res, fmt.Sprintf("[[%v] %v/](fg:green)", counter, file.Name()))
			}
		}
		for _, file := range files {
			if strings.HasPrefix(file.Name(), ".") {
				continue
			}
			if !file.IsDir() {
				counter++
				res = append(res, fmt.Sprintf("[%v] %v", counter, file.Name()))
			}
		}
	}

	return res
}
