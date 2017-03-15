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
	Versions []version
}

type version struct {
	Identifier         string
	WithPersistentData bool
	WithCode           bool
}

func (app Application) GetLatestVersion() version {
	return app.Versions[len(app.Versions)-1]
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

func convertItemsToVersions(items []stow.Item) (versions []version) {
	versionItemBuckets := mapItemsByVersion(items)

	versionIdentifiers := make([]string, len(versionItemBuckets))
	i := 0
	for k := range versionItemBuckets {
		versionIdentifiers[i] = k
		i++
	}
	sort.Strings(versionIdentifiers)

	for i := 0; i < len(versionIdentifiers); i++ {
		versionIdentifier := versionIdentifiers[i]
		items := versionItemBuckets[versionIdentifier]

		hasPersistentData := itemsContainItemWithSuffix(items, "-persistent-data.tar.gz.gpg")
		hasCode := itemsContainItemWithSuffix(items, "-code.tar.gz.gpg")

		versions = append(versions, version{
			Identifier: versionIdentifier,
			WithPersistentData: hasPersistentData,
			WithCode: hasCode,
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
	versionPattern := regexp.MustCompile(`(.*)__.*?__.*?(-manifest\.json\.gpg|-persistent-data\.tar\.gz|-code\.tar\.gz)`)

	return versionPattern.FindStringSubmatch(item.Name())[1]
}
func getVersionIdentifier(item stow.Item) string {
	versionPattern := regexp.MustCompile(`(.*)(-manifest\.json\.gpg|-persistent-data\.tar\.gz|-code\.tar\.gz)`)

	return versionPattern.FindStringSubmatch(item.Name())[1]
}
func itemsContainItemWithSuffix(items []stow.Item, suffix string) bool {
	for _, item := range items {
		if strings.HasSuffix(item.Name(), suffix) {
			return true
		}
	}

	return false
}