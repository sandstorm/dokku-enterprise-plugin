package cloud

import (
	"github.com/sandstorm/dokku-enterprise-plugin/core/cloudStorage"
	"log"
	"github.com/sandstorm/dokku-enterprise-plugin/core/utility"
	"strconv"
)

func List() {
	allCloudApps, err := cloudStorage.GetAllApplications()
	if err != nil {
		log.Fatal(err)
	}

	headerRow := []string{"NAME", "VERSIONS", "LATEST"}

	contentRows := make([][]string, len(allCloudApps))
	for i := 0; i < len(contentRows); i++ {
		contentRows[i] = convertApplicationToTableRow(allCloudApps[i])
	}

	utility.RenderAsTable(headerRow, contentRows)
}

func convertApplicationToTableRow(application cloudStorage.Application) (result []string) {
	result = append(result, application.Name)
	result = append(result, strconv.Itoa(len(application.Versions)))
	result = append(result, application.GetLatestVersion().Identifier)

	return
}