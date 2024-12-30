/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path"
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

	//rootCmd.AddCommand(updateLocalConfig)

	updateLocalConfig.Flags().StringVar(&enumValue, "action", "", fmt.Sprintf("Specify the action (%v)", allowedValuesInLocalUpdate))
	updateLocalConfig.MarkFlagRequired("action")

}

func cleanupDir(dir string) {
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

		os.RemoveAll(fileSystemPath)
	}

}

func cleanDir(dir string) {
	dirsToClean := findAllChildDirs(dir)

	dirsToCleanCount := len(dirsToClean)

	if dirsToCleanCount <= 0 {
		printError("Nothing to clean...")
	}

	bar := GenerateProgressBar(dirsToCleanCount, "Cleaning Dirs")
	for i := 0; i < dirsToCleanCount; i++ {
		dir := dirsToClean[i]
		cleanupDir(dir)
		bar.Add(i + 1)
	}
}

func updateLocalConfigurationWithRemoteConfig() {
	configFilePath := generateConfigPaths()
	downloadConfigFile(config.RemoteConfigURL, configFilePath)
}

func updateValuesInLocalConfig() {

	if enumValue == "add" {
	} else if enumValue == "update" {
	} else if enumValue == "remove" {
	}

}
