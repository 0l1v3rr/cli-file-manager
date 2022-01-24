package pkg

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

func NoEx(p string) []string {
	files, err := ioutil.ReadDir(p)
	if err != nil {
		log.Fatal(err)
	}

	res := []string{}
	counter := 0

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
			if file.Name()[0] == '.' {
				res = append(res, fmt.Sprintf("[%v] %v", counter, file.Name()))
			} else if !strings.Contains(file.Name(), ".") {
				res = append(res, fmt.Sprintf("[%v] %v", counter, file.Name()))
			} else {
				name := strings.Split(file.Name(), ".")
				name = name[:len(name)-1]
				res = append(res, fmt.Sprintf("[%v] %v", counter, strings.Join(name, ".")))
			}
		}
	}

	return res
}

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
