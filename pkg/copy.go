package pkg

import (
	"io/ioutil"
	"strings"
)

func Copy(src string, path string) error {
	input, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path+strings.Split(src, "/")[len(strings.Split(src, "/"))-1], input, 0644)
	if err != nil {
		return err
	}

	return nil
}
