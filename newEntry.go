package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-yaml/yaml"
)

const template string = `filename: '' # short description not used in final changlog
entrytype: # one of either (added, changed, deprecated, removed, fixed, security)
change: '' # the actual change entry, must be in the quotes


# :!:  Changelog Entry Creator  :!:
# ---------------------------------
#
#
# filetype:
#   - a short name for the file to be created
# entrytype:
#   - either added, changed, deprecated, removed, fixed, security
# change:
#   - the actual line you want to put in the changelog
#
# Example:
#
# filename: attachment type field
# entrytype: added
# change: Attachment type field added to the Attachment Table
`

type ChangelogEntry struct {
	Filename  string
	EntryType string
	Change    string
}

func CreateChangelogEntry(root string) error {
	tmpFile, err := os.CreateTemp(os.TempDir(), "changelog*.yml")
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %v", err)
	}

	tmpFile.WriteString(template)
	tmpFile.Close()

	editor := os.ExpandEnv("$EDITOR")
	if editor == "" {
		editor = "nano"
	}

	cmd := exec.Command(editor, tmpFile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run exec command: %v", err)
	}

	tmpFile, err = os.Open(tmpFile.Name())
	if err != nil {
		return fmt.Errorf("failed to open temporary file: %v", err)
	}

	var entry ChangelogEntry
	if err := yaml.NewDecoder(tmpFile).Decode(&entry); err != nil {
		return fmt.Errorf("invalid YAML: %e", err)
	}

	// make a folder if it doesn't exist
	changelogDir := filepath.Join(root, ".changelog")
	unreleasedDir := filepath.Join(changelogDir, "Unreleased")
	os.MkdirAll(unreleasedDir, os.ModePerm)

	filename := strings.ReplaceAll(entry.Filename, " ", "-")
	entryType := strings.ToLower(entry.EntryType)

	newFileName := filepath.Join(unreleasedDir, fmt.Sprintf("%s-%s.yml", time.Now().Format("2006-01-02"), filename))
	newFile, err := os.Create(newFileName)
	if err != nil {
		return fmt.Errorf("failed to create changelog file: %v", err)
	}

	newFile.WriteString(fmt.Sprintf("%s:\n- '%s'", entryType, escapeYaml(entry.Change)))
	newFile.Close()

	fmt.Printf("%s created successfully!\n", newFile.Name())

	return nil
}

func escapeYaml(s string) string {
	return strings.ReplaceAll(s, "'", `\'`)
}
