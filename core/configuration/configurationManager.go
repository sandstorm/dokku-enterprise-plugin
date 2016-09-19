package configuration

import (
	"encoding/json"
	"github.com/kardianos/osext"
	"io/ioutil"
)

var configurationCache configuration = configuration{
	ApiEndpointUrl: "",
}

var isInitialized bool = false

func Get() configuration {

	if !isInitialized {
		executableFolder, _ := osext.ExecutableFolder()
		configBytes, _ := ioutil.ReadFile(executableFolder + "/config.json")
		json.Unmarshal(configBytes, &configurationCache)
		isInitialized = true
	}

	return configurationCache
}
