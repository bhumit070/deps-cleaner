package cmd

var dependenciesFileURL string = "https://raw.githubusercontent.com/bhumit070/deps-cleaner/refs/heads/main/cmd/config/config.json"

var skipConfirmation bool = false
var updateConfig string

type Config struct {
	Deps   map[string]bool `json:"deps"`
	Ignore map[string]bool `json:"ignore"`
}

var config Config
