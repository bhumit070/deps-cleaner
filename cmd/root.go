package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/fatih/color"
	"github.com/rodaine/table"
)

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	loadConfigFiles()

	rootCmd.Flags().BoolP("version", "v", false, "Prints the version")

	rootCmd.AddCommand(updateLocalConfigurationWithRemote)

	rootCmd.AddCommand(cleanDependencies)
	cleanDependencies.Flags().BoolVarP(&skipConfirmation, "yes", "y", false, "Do not ask for confirmation before removing the dependencies")
	cleanDependencies.Flags().BoolVarP(&isDryRun, "dry-run", "", false, "Prints all the file that will be deleted")

	rootCmd.AddCommand(updateLocalConfig)

	updateLocalConfig.Flags().StringVar(&enumValue, "action", "", fmt.Sprintf("Specify the action (%v)", allowedValuesInLocalUpdate))
	updateLocalConfig.MarkFlagRequired("action")

}

func cleanupDir(dir string, keepTheFile bool, filesToRemove *[]string) {
	dir, isValid := validatePath(dir, false, true)

	if !isValid {
		printError("directory " + dir + " is invalid.")
	}

	for dependencyFolder := range config.Deps {
		folderPath := path.Join(dir, dependencyFolder)

		fileSystemPath, isValid := validatePath(folderPath, true, false)

		if !isValid {
			continue
		}

		if !keepTheFile {

			if size, err := getDirSize(fileSystemPath); err == nil {
				totalFreedSpace += size
			}

			os.RemoveAll(fileSystemPath)
		}

		*filesToRemove = append(*filesToRemove, fileSystemPath)

	}
}

func cleanDir(dir string) map[string]string {
	totalFreedSpace = 0
	cleanedDirs := map[string]string{}
	dirsToClean := findAllChildDirs(dir)
	dirsToCleanCount := len(dirsToClean)

	if dirsToCleanCount <= 0 {
		printError("Nothing to clean...")
	}

	message := "Cleaning Dirs..."

	if isDryRun {
		message = "Scanning dirs..."
	}

	var filesToRemove []string

	bar := GenerateProgressBar(dirsToCleanCount, message)
	for i := 0; i < dirsToCleanCount; i++ {
		dir := dirsToClean[i]

		cleanupDir(dir, isDryRun, &filesToRemove)

		bar.Add(i + 1)
	}
	bar.Close()
	fmt.Println()

	if len(filesToRemove) <= 0 {
		printError("Nothing to clean")
	}

	if isDryRun {

		headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
		columnFmt := color.New(color.FgYellow).SprintfFunc()

		tbl := table.New("FilePath")
		tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

		for i := range filesToRemove {
			tbl.AddRow(filesToRemove[i])
		}

		fmt.Println()
		tbl.Print()
		fmt.Println()

		promptForConfirmation("Above directories will be removed.")

		for i := range filesToRemove {

			if size, err := getDirSize(filesToRemove[i]); err == nil {
				totalFreedSpace += size
			}

			os.RemoveAll(filesToRemove[i])
		}

	}

	fmt.Printf("\nTotal space freed: %s\n", formatBytes(totalFreedSpace))

	return cleanedDirs
}

func updateLocalConfigurationWithRemoteConfig() {
	configFilePath := generateConfigPaths()
	downloadConfigFile(config.RemoteConfigURL, configFilePath)
}

func updateValuesInLocalConfig(args []string) {

	var configKey, key string

	if len(args) >= 1 {
		configKey = args[0]
	}

	if len(args) >= 2 {
		key = args[1]
	}

	isConfigValueChanged := false

	if configKey == "" || key == "" {
		printError("values are missing to add value to config")
	}

	if enumValue == "add" {

		if configKey == "deps" {
			config.Deps[key] = true
			isConfigValueChanged = true
		}

		if configKey == "ignore" {
			config.Ignore[key] = true
			isConfigValueChanged = true
		}

	} else if enumValue == "remove" {

		if configKey == "deps" {
			delete(config.Deps, key)
			isConfigValueChanged = true
		}

		if configKey == "ignore" {
			delete(config.Ignore, key)
			isConfigValueChanged = true
		}
	}

	if !isConfigValueChanged {
		return
	}

	configFilePath := generateConfigPaths()

	outFile, err := os.Create(configFilePath)
	if err != nil {
		printError("Failed to load config file")
	}
	defer outFile.Close()

	encoder := json.NewEncoder(outFile)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(config); err != nil {
		printError("Failed to save the file.")
	}

}
