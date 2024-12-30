/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"
	"path"
)

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
	for i := 0; i < dirsToCleanCount; i++ {
		dir := dirsToClean[i]
		cleanupDir(dir)
	}
}

func updateLocalConfigurationWithRemoteConfig() {
	configFilePath := generateConfigPaths()
	downloadConfigFile(dependenciesFileURL, configFilePath)
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("version", "v", false, "Prints the version")
	rootCmd.Flags().BoolVarP(&skipConfirmation, "yes", "y", false, "Do not ask for confirmation before removing the dependencies")

	loadConfigFiles()

	rootCmd.AddCommand(updateLocalConfigurationWithRemote)
	rootCmd.AddCommand(cleanDependencies)

}
