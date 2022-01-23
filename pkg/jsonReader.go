package pkg

import (
	"fmt"
	"io/ioutil"

	"github.com/Jeffail/gabs"
)

func ReadJson() (string, error) {
	data, err1 := ioutil.ReadFile("/usr/local/cli-file-manager/settings.json")
	if err1 != nil {
		return "", err1
	}

	jsonParsed, err2 := gabs.ParseJSON([]byte(data))
	if err2 != nil {
		return "", err2
	}

	return fmt.Sprintf("%v", jsonParsed.Path("defaultPanel").Data()), err2
}
