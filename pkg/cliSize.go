package pkg

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func getCliSize() []string {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}

	s := string(out)
	s = strings.TrimSpace(s)
	return strings.Split(s, " ")
}

func GetCliHeight() int {
	res, err := strconv.Atoi(getCliSize()[0])
	if err != nil {
		fmt.Println(err)
	}
	return res
}

func GetCliWidth() int {
	res, err := strconv.Atoi(getCliSize()[1])
	if err != nil {
		fmt.Println(err)
	}
	return res
}
