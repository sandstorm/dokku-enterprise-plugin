package manifest

type manifest struct {
	Mariadb []string    `json:"mariadb,omitempty"`
}

type manifestWrapper struct {
	Version   int `json:"version"`
	AppName   string `json:"appName"`
	Manifest  *manifest `json:"manifest"`
	Errors    []string `json:"errors,omitempty"`
	DebugInfo []string `json:"debugInfo,omitempty"`
}