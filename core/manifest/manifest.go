package manifest

type manifestWrapper struct {
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