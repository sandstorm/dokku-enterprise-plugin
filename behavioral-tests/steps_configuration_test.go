package main

import (
	"encoding/json"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/kardianos/osext"
	"io/ioutil"
)

// create a configuration JSON from a data table
func theConfigurationIs(configuration *gherkin.DataTable) error {
	config := make(map[string]string)

	for _, row := range configuration.Rows {
		config[row.Cells[0].Value] = row.Cells[1].Value
	}

	configAsBytes, _ := json.MarshalIndent(config, "", "    ")

	behavioralTestsFolder, _ := osext.ExecutableFolder()
	ioutil.WriteFile(behavioralTestsFolder+"/../bin-build/config.json", configAsBytes, 0644)

	return nil
}


func theCloudConfigurationIs(cloudConfiguration *gherkin.DataTable) error {
	cloudConfig := make(map[string]string)

	for _, row := range cloudConfiguration.Rows {
		cloudConfig[row.Cells[0].Value] = row.Cells[1].Value
	}

	cfg := configWrapper {
		CloudBackup: cloudConfig,
	}

	configAsBytes, _ := json.MarshalIndent(cfg, "", "    ")

	behavioralTestsFolder, _ := osext.ExecutableFolder()
	ioutil.WriteFile(behavioralTestsFolder+"/../bin-build/config.json", configAsBytes, 0644)

	return nil
}

type configWrapper struct {
	CloudBackup    map[string]string `json:"cloudBackup"`
}