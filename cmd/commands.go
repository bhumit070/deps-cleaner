package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "deps-cleaner",
	Short: "Clean your space by removing dependencies folders.",
	Long:  `deps-cleaner is a cli that helps you in cleaning your space by removing dependencies folders that can be recreated very easily.`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) <= 0 {
			cmd.Help()
			return
		}

	},
	Version: "0.0.1",
}

var updateLocalConfigurationWithRemote = &cobra.Command{
	Use:   "update:config",
	Short: "Update local configuration with remote configuration",
	Run: func(cmd *cobra.Command, args []string) {
		updateLocalConfigurationWithRemoteConfig()
		fmt.Println("Config updated successfully...")
	},
}

var cleanDependencies = &cobra.Command{
	Use:   "clean",
	Short: "Clean dependencies",
	Run: func(cmd *cobra.Command, args []string) {
		if updateConfig != "" {
			loadConfigFiles()
		}

		dirToClean := "."
		if len(args) > 0 {
			dirToClean = args[len(args)-1]
		}

		dirToClean, isInvalid := validatePath(dirToClean, false, true)

		if !isInvalid {
			printError("Directory " + dirToClean + " is invalid.")
		}

		if skipConfirmation == false {
			promptForConfirmation(dirToClean)
		}

		cleanDir(dirToClean)
	},
}

var updateLocalConfig = &cobra.Command{
	Use:   "config",
	Short: "Update local configuration with your custom inputs",
	Run: func(cmd *cobra.Command, args []string) {
		updateValuesInLocalConfig()
	},
}
