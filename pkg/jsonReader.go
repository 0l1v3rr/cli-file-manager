package pkg

import (
	"fmt"
	"io/ioutil"

	"github.com/Jeffail/gabs"
)

func ReadJson() (string, error) {
	data, err1 := ioutil.ReadFile("./settings.json")
	if err1 != nil {
		fmt.Print(err1)
	} else {
		jsonParsed, err2 := gabs.ParseJSON([]byte(data))
		if err2 != nil {
			panic(err2)
		}
		return fmt.Sprintf("%v", jsonParsed.Path("defaultPanel").Data()), err2
	}

	return "", err1
}
