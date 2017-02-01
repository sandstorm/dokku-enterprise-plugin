package cloud

import (
	"log"
	"github.com/graymeta/stow"
	"strings"
	"regexp"
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

func List() {
	//todo
}

func getApplication(applicationName string, container stow.Container) Application {
	var limit int = 1000 // todo: should we support "batch processing" of items?

	items, _, err := container.Items(applicationName, stow.CursorStart, limit)
	if err != nil {
		log.Fatalf("ERROR: could not access items of Cloud Storage, error was: %v", err)
	}

	return Application{Name: applicationName, Versions: convertItemsToVersions(items)}
}
func convertItemsToVersions(items []stow.Item) (versions []version) {
	versionItemBuckets := mapItemsByVersion(items)

	for versionIdentifier, items := range versionItemBuckets {
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
func mapItemsByVersion(items []stow.Item) (itemsByVersions map[string][]stow.Item) {
	itemsByVersions = make(map[string][]stow.Item)

	for _, item := range items {
		version := getVersionIdentifier(item)
		if itemsByVersions[version] == nil {
			itemsByVersions[version] = make([]stow.Item, 0)
		}
		itemsByVersions[version] = append(itemsByVersions[version], item)
	}

	return
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