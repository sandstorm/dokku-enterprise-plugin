package cloud

import (
	"github.com/sandstorm/dokku-enterprise-plugin/core/cloudStorage"
	"log"
	"github.com/sandstorm/dokku-enterprise-plugin/core/utility"
	"strconv"
)

func Info(applicationName string) {
	application, err := cloudStorage.GetApplication(applicationName)
	if err != nil {
		log.Fatal(err)
	}

	headerRow := []string{"VERSION", "DATE", "CODE SIZE (KB)", "DATA SIZE (KB)"}

	contentRows := make([][]string, len(application.Versions))
	for i := 0; i < len(contentRows); i++ {
		contentRows[i] = convertVersionToTableRow(application.Versions[i])
	}

	utility.RenderAsTable(headerRow, contentRows)
}

func convertVersionToTableRow(version cloudStorage.Version) (result []string) {
	result = append(result, version.Identifier)
	result = append(result, version.GetDate())
	result = append(result, strconv.FormatFloat(version.CodeSizeInKb, 'f', 3, 64))
	result = append(result, strconv.FormatFloat(version.DataSizeInKb, 'f', 3, 64))

	return
}