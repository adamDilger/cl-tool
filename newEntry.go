package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-yaml/yaml"
)

const template string = `# :!:  Changelog Entry Creator  :!:
# ---------------------------------
#
# Fill out the fields below, then save and quit

filename: ''
entrytype:
change: ''


# Help: --------------------------------
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
# filename: 'attachment type field'
# entrytype: 'added'
# change: 'Attachment type field added to the Attachment Table'
`

var allEntryTypes []string = []string{
	"added",
	"changed",
	"deprecated",
	"removed",
	"fixed",
	"security",
}

type ChangelogEntry struct {
	Filename  string
	EntryType string
	Change    string
}

func (ce *ChangelogEntry) IsValid() error {
	if ce.Filename == "" {
		return errors.New("filename must not be empty")
	}

	if ce.Change == "" {
		return errors.New("change must not be empty")
	}

	found := false
	for _, e := range allEntryTypes {
		if e == ce.EntryType {
			found = true
		}
	}

	if !found {
		return fmt.Errorf("invalid entry type: %s", ce.EntryType)
	}

	return nil
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

	repeat := true
	var entry ChangelogEntry

	for repeat {
		repeat = false

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

		if err := yaml.NewDecoder(tmpFile).Decode(&entry); err != nil {
			fmt.Printf("Invalid YAML: %v\n", err)
			repeat = userYesOrNo("Do you want to retry?")
		}

		if err := entry.IsValid(); err != nil {
			fmt.Printf("Invalid data: %v\n", err)
			repeat = userYesOrNo("Do you want to retry?")
		}
	}

	// make a folder if it doesn't exist
	changelogDir := filepath.Join(root, ".changelog")
	unreleasedDir := filepath.Join(changelogDir, "Unreleased")
	os.MkdirAll(unreleasedDir, os.ModePerm)

	filename := strings.ReplaceAll(entry.Filename, " ", "-")
	entryType := strings.ToLower(string(entry.EntryType))

	newFileName := filepath.Join(unreleasedDir, fmt.Sprintf("%s-%s.yml", time.Now().Format("2006-01-02"), filename))
	newFile, err := os.Create(newFileName)
	if err != nil {
		return fmt.Errorf("failed to create changelog file: %v", err)
	}

	newFile.WriteString(fmt.Sprintf("%s:\n  - '%s'\n", entryType, escapeYaml(entry.Change)))
	newFile.Close()

	fmt.Printf("%s created successfully!\n", newFile.Name())

	return nil
}

func escapeYaml(s string) string {
	return strings.ReplaceAll(s, "'", `\'`)
}

func userYesOrNo(question string) bool {
	fmt.Printf("%s (y/n): ", question)
	scanner := bufio.NewScanner(os.Stdin)

	if scanner.Scan() {
		text := scanner.Text()
		return text == "y" || text == "Y"
	}

	return false
}
