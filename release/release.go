package release

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

func CreateRelease(root, version string) error {
	var err error

	versionNumber := version
	if versionNumber == "" {
		versionNumber, err = getVersionNumber()
		if err != nil {
			return err
		}
	} else {
		if err := validateVersionNumber(versionNumber); err != nil {
			return err
		}
	}

	fmt.Printf("Renaming Unreleased folder to %s - %s\n", versionNumber, time.Now().Format("2006-01-02"))
	err = renameUnreleasedFolder(root, versionNumber, time.Now().Format("2006-01-02"))
	if err != nil {
		return err
	}

	fmt.Println("Success!")
	return nil
}

func getVersionNumber() (string, error) {
	fmt.Print("Version Number: ")

	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		return "", fmt.Errorf("error parsing input: %v", scanner.Err())
	}

	versionNumber := scanner.Text()

	if err := validateVersionNumber(versionNumber); err != nil {
		return versionNumber, err
	}

	return versionNumber, nil
}

func renameUnreleasedFolder(path, version, date string) error {
	og := filepath.Join(path, ".changelog", "Unreleased")
	new := filepath.Join(path, ".changelog", fmt.Sprintf("%s_%s", version, date))
	return os.Rename(og, new)
}

func validateVersionNumber(versionNumber string) error {
	if !regexp.MustCompile(`^\d+\.\d+\.\d+$`).MatchString(versionNumber) {
		return fmt.Errorf("invalid version number: %s", versionNumber)
	}

	return nil
}
