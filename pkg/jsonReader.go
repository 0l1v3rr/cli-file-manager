package pkg

import (
	"fmt"
	"io/ioutil"

	"github.com/Jeffail/gabs"
)

func ReadJson() string {
	data, err := ioutil.ReadFile("./settings.json")
	if err != nil {
		fmt.Print(err)
	}

	jsonParsed, err := gabs.ParseJSON([]byte(data))
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("%v", jsonParsed.Path("defaultPanel").Data())
}
