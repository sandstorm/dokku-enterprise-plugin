package main

import (
	"github.com/DATA-DOG/godog/gherkin"
	"encoding/json"
	"io/ioutil"
	"github.com/kardianos/osext"
)

func theConfigurationIs(configuration *gherkin.DataTable) error {
	config := make(map[string]string)

	for _, row := range configuration.Rows {
		config[row.Cells[0].Value] = row.Cells[1].Value
	}

	configAsBytes, _ := json.MarshalIndent(config, "", "    ")

	behavioralTestsFolder, _ := osext.ExecutableFolder()
	ioutil.WriteFile(behavioralTestsFolder + "/../bin-build/config.json", configAsBytes, 0644)

	return nil
}