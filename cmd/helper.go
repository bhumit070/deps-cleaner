package cmd

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/schollz/progressbar/v3"
)

func findAllChildDirs(rootDir string) []string {
	var dirs []string

	if isReadAndWriteAccessGranted := checkReadAndWriteAccess(rootDir); !isReadAndWriteAccessGranted {
		printError("Read & Write access is not available at " + rootDir)
	}

	err := filepath.Walk(rootDir, func(fsPath string, info os.FileInfo, err error) error {

		if info == nil {
			return nil
		}

		name := info.Name()
		if config.Ignore[name] || config.Deps[name] || config.Ignore[fsPath] {
			return fs.SkipDir
		}

		if isReadAndWriteAccessGranted := checkReadAndWriteAccess(fsPath); !isReadAndWriteAccessGranted {
			return nil
		}

		if err != nil {
			return err
		}

		if info.IsDir() {
			dirs = append(dirs, fsPath)
		}
		return nil
	})

	if err != nil {
		printError(err.Error())
	}

	return dirs
}

func validatePath(filePath string, checkExistenceOnly bool, checkDirOnly bool) (string, bool) {

	if checkDirOnly {
		if filePath == "." {
			dirPath, errInGettingCWD := os.Getwd()

			if errInGettingCWD != nil {
				printError(errInGettingCWD.Error())
			}

			filePath = dirPath
		}
	}

	stats, errorInStats := os.Stat(filePath)

	if errorInStats != nil {

		if checkExistenceOnly {
			if os.IsNotExist(errorInStats) {
				return "", false
			}
		}

		printError(errorInStats.Error())
	}

	if checkDirOnly {

		if stats.IsDir() == false {
			return filePath, true
		}
	}

	return filePath, true
}

func printError(message string) {
	fmt.Println(message)
	os.Exit(0)
}

func promptForConfirmation(label string) {

	label += " (y/N): "

	prompt := promptui.Prompt{
		Label:     label,
		Default:   "n",
		IsConfirm: true,
	}

	result, err := prompt.Run()

	result = strings.ToLower(result)

	if err != nil {
		printError(err.Error())
	}

	if result == "n" {
		os.Exit(0)
	}
}

func generateConfigPaths() string {
	userHomeDir, errorInGettingUserHomeDir := os.UserHomeDir()

	if errorInGettingUserHomeDir != nil {
		printError(errorInGettingUserHomeDir.Error())
	}

	configFolder := path.Join(userHomeDir, ".config", "deps-cleaner")

	_, isValid := validatePath(configFolder, true, true)

	if !isValid {
		errorInCreatingConfigFolder := os.MkdirAll(configFolder, os.ModePerm)
		if errorInCreatingConfigFolder != nil {
			printError(errorInGettingUserHomeDir.Error())
		}
		isValid = true
	}

	configFilePath := path.Join(configFolder, "config.json")

	return configFilePath
}

func loadConfigFiles() {

	configFilePath := generateConfigPaths()

	_, isValid := validatePath(configFilePath, true, false)

	if !isValid {
		downloadConfigFile(remoteConfigURL, configFilePath)
	}

	file, err := os.Open(configFilePath)
	if err != nil {
		printError(err.Error())
		return
	}

	defer file.Close()

	decoder := json.NewDecoder(file)

	if err := decoder.Decode(&config); err != nil {
		printError(err.Error())
	}

	if config.RemoteConfigURL == "" {
		config.RemoteConfigURL = remoteConfigURL
	}

}

func downloadConfigFile(url, filepath string) error {
	existingData := make(map[string]interface{})
	if _, err := os.Stat(filepath); err == nil {
		file, err := os.Open(filepath)
		if err != nil {
			return fmt.Errorf("failed to open existing file: %v", err)
		}
		defer file.Close()

		decoder := json.NewDecoder(file)
		if err := decoder.Decode(&existingData); err != nil {
			return fmt.Errorf("failed to decode existing JSON: %v", err)
		}
	}

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to make GET request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download file: status code %d", resp.StatusCode)
	}

	newData := make(map[string]interface{})
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&newData); err != nil {
		return fmt.Errorf("failed to decode new JSON: %v", err)
	}

	mergeMaps(existingData, newData)

	if _, exists := existingData["remoteConfigURL"]; !exists {
		existingData["remoteConfigURL"] = remoteConfigURL
	}

	outFile, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer outFile.Close()

	encoder := json.NewEncoder(outFile)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(existingData); err != nil {
		return fmt.Errorf("failed to write merged JSON to file: %v", err)
	}

	return nil
}

func mergeMaps(dest, src map[string]interface{}) {
	for key, srcValue := range src {
		if srcMap, ok := srcValue.(map[string]interface{}); ok {
			if destValue, ok := dest[key].(map[string]interface{}); ok {
				mergeMaps(destValue, srcMap)
			} else {
				dest[key] = srcMap
			}
		} else {
			dest[key] = srcValue
		}
	}
}

func GenerateProgressBar(barPercentage int, message string) *progressbar.ProgressBar {

	bar := progressbar.NewOptions(barPercentage,
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetWidth(30),
		progressbar.OptionSetDescription(message),
		progressbar.OptionSetPredictTime(false),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))

	return bar
}

func getDirSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	return size, err
}

func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func checkReadAndWriteAccess(folderPath string) bool {
	// Check Read Access
	f, err := os.Open(folderPath)
	if err != nil {
		return false
	}
	f.Close() // Close immediately

	// Check Write Access
	tempTextFile := path.Join(folderPath, "temp.txt")
	f, err = os.Create(tempTextFile)
	if err != nil {
		return false
	}
	f.Close() // Close immediately

	// Cleanup
	err = os.Remove(tempTextFile)
	return err == nil
}
