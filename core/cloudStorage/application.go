package cloudStorage

import (
	"github.com/graymeta/stow"
	"regexp"
	"strings"
	"fmt"
	"sort"
)

type Application struct {
	Name     string
	Versions []Version
}

type Version struct {
	Identifier	string
	CodeSizeInKb	float64
	DataSizeInKb    float64
}

func (version Version) GetDate() string {
	datePattern := regexp.MustCompile(`.*__(.*?)__.*?`)
	return datePattern.FindStringSubmatch(version.Identifier)[1]
}

const ITEM_LIMIT = 1000 // todo: should we support "batch processing" of items?

func GetAllApplications() (result []Application, err error) {
	items, err := retrieveItemsFromCloudStorage("")
	if err != nil {
		return nil, fmt.Errorf("Could not retrieve any applications from cloud storage, error was: %v", err)
	}

	itemsByApplicationName := mapItemsByApplicationName(items)

	applicationNames := make([]string, len(itemsByApplicationName))
	i := 0
	for k := range itemsByApplicationName {
		applicationNames[i] = k
		i++
	}
	sort.Strings(applicationNames)

	for i := 0; i < len(applicationNames); i++ {
		applicationName := applicationNames[i]
		items = itemsByApplicationName[applicationName]

		result = append(result, Application{Name: applicationName, Versions: convertItemsToVersions(items)})
	}

	return
}

func GetApplication(applicationName string) (Application, error) {
	var app Application

	items, err := retrieveItemsFromCloudStorage(applicationName)
	if err != nil {
		return app, fmt.Errorf("Could not retrieve application from cloud storage, error was: %v", err)
	}

	app = Application{Name: applicationName, Versions: convertItemsToVersions(items)}
	return app, nil
}

func retrieveItemsFromCloudStorage(prefix string) ([]stow.Item, error) {
	items, _, err := getContainer().Items(prefix, stow.CursorStart, ITEM_LIMIT)
	return items, err
}

func convertItemsToVersions(items []stow.Item) (versions []Version) {
	versionItemBuckets := mapItemsByVersion(items)

	versionIdentifiers := make([]string, len(versionItemBuckets))
	i := 0
	for k := range versionItemBuckets {
		versionIdentifiers[i] = k
		i++
	}
	sort.Sort(sort.Reverse(sort.StringSlice(versionIdentifiers)))

	for i := 0; i < len(versionIdentifiers); i++ {
		versionIdentifier := versionIdentifiers[i]
		items := versionItemBuckets[versionIdentifier]

		codeFile := getItemWithSuffix(items, "-persistent-data.tar.gz.gpg")
		var codeSizeInKb float64 = 0
		if codeFile != nil {
			size, err := codeFile.Size()
			if err == nil {
				codeSizeInKb = float64(size) / 1024
			}
		}

		dataFile := getItemWithSuffix(items, "-code.tar.gz.gpg")
		var dataSizeInKb float64 = 0
		if dataFile != nil {
			size, err := dataFile.Size()
			if err == nil {
				dataSizeInKb = float64(size) / 1024
			}
		}

		versions = append(versions, Version{
			Identifier: versionIdentifier,
			CodeSizeInKb: codeSizeInKb,
			DataSizeInKb: dataSizeInKb,
		})
	}

	return
}
func mapItemsByApplicationName(items []stow.Item) (itemsByApplicationName map[string][]stow.Item) {
	itemsByApplicationName = make(map[string][]stow.Item)

	for _, item := range items {
		applicationName := getApplicationName(item)
		if itemsByApplicationName[applicationName] == nil {
			itemsByApplicationName[applicationName] = make([]stow.Item, 0)
		}
		itemsByApplicationName[applicationName] = append(itemsByApplicationName[applicationName], item)
	}

	return
}
func mapItemsByVersion(items []stow.Item) (itemsByVersion map[string][]stow.Item) {
	itemsByVersion = make(map[string][]stow.Item)

	for _, item := range items {
		version := getVersionIdentifier(item)
		if itemsByVersion[version] == nil {
			itemsByVersion[version] = make([]stow.Item, 0)
		}
		itemsByVersion[version] = append(itemsByVersion[version], item)
	}

	return
}
func getApplicationName(item stow.Item) string {
	// version-name pattern: <appName>__<date_time>__<ip>-<type>
	applicationNamePattern := regexp.MustCompile(`(.*)__.*?__.*?(-manifest\.json\.gpg|-persistent-data\.tar\.gz|-code\.tar\.gz)`)
	return applicationNamePattern.FindStringSubmatch(item.Name())[1]
}
func getVersionIdentifier(item stow.Item) string {
	versionPattern := regexp.MustCompile(`(.*)(-manifest\.json\.gpg|-persistent-data\.tar\.gz|-code\.tar\.gz)`)
	return versionPattern.FindStringSubmatch(item.Name())[1]
}
func getItemWithSuffix(items []stow.Item, suffix string) stow.Item {
	for _, item := range items {
		if strings.HasSuffix(item.Name(), suffix) {
			return item
		}
	}
	return nil
}