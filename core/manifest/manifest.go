package manifest

import (
	"strings"
	"encoding/json"
	"log"
)

type ManifestWrapper struct {
	Version   int `json:"version"`
	AppName   string `json:"appName"`
	Manifest  *manifest `json:"manifest"`
	Errors    []string `json:"errors,omitempty"`
	DebugInfo []string `json:"debugInfo,omitempty"`
}

type manifest struct {
	Mariadb       []string    `json:"mariadb,omitempty"`
	DockerOptions dockerOptions `json:"dockerOptions,omitempty"`
	Config        map[string]string `json:"config,omitempty"`
}

type dockerOptions struct {
	Build  []string     `json:"build,omitempty"`
	Deploy []string     `json:"deploy,omitempty"`
	Run    []string     `json:"run,omitempty"`
}

func ReplaceAppNamePlaceholder(input, applicationName string) string {
	return strings.Replace(input, "[appName]", applicationName, -1)
}

func SerializeManifest(manifestWrapper ManifestWrapper) []byte {
	manifestAsBytes, err := json.MarshalIndent(manifestWrapper, "", "  ")

	if err != nil {
		log.Fatalf("There was an error serializing JSON manifest: %v", err)
	}

	return manifestAsBytes
}

func DeserializeManifest(manifestAsBytes []byte) ManifestWrapper {
	manifestWrapper := ManifestWrapper{}

	err := json.Unmarshal(manifestAsBytes, &manifestWrapper)

	if err != nil {
		log.Fatal("ERROR: JSON could not be parsed")
	}

	return manifestWrapper
}

func (manifestWrapper ManifestWrapper) GetDockerOptions() []string {
	allDockerOptions := make([]string, 0, 20)
	allDockerOptions = append(allDockerOptions, manifestWrapper.Manifest.DockerOptions.Build...)
	allDockerOptions = append(allDockerOptions, manifestWrapper.Manifest.DockerOptions.Deploy...)
	allDockerOptions = append(allDockerOptions, manifestWrapper.Manifest.DockerOptions.Run...)
	return removeDuplicates(allDockerOptions)
}

func removeDuplicates(sliceWithDuplicates []string) []string {
	found := make(map[string]bool)
	sliceWithoutDuplicates := make([]string, len(sliceWithDuplicates))

	j := 0
	for i, x := range sliceWithDuplicates {
		if !found[x] {
			found[x] = true
			sliceWithoutDuplicates[j] = sliceWithDuplicates[i]
			j++
		}
	}

	return sliceWithoutDuplicates[:j]
}