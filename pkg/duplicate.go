package pkg

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func Duplicate(src string, path string) error {
	input, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}

	filename := strings.Split(src, "/")[len(strings.Split(src, "/"))-1]
	if filename[0] == '.' {
		err = ioutil.WriteFile(path+"/"+fmt.Sprintf("%s-copy", filename), input, 0644)
		if err != nil {
			return err
		}
	} else if !strings.Contains(filename, ".") {
		err = ioutil.WriteFile(path+"/"+fmt.Sprintf("%s-copy", filename), input, 0644)
		if err != nil {
			return err
		}
	} else {
		split := strings.Split(filename, ".")
		ex := split[len(split)-1]
		split = split[:len(split)-1]
		err = ioutil.WriteFile(path+"/"+fmt.Sprintf("%s-copy.%s", strings.Join(split, "."), ex), input, 0644)
		if err != nil {
			return err
		}
	}

	return nil
}
