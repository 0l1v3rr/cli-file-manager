package pkg

import (
	"fmt"
	"io/ioutil"
	"log"
)

func ReadFiles(p string) []string {
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
			res = append(res, fmt.Sprintf("[%v] %v", counter, file.Name()))
		}
	}

	return res
}
