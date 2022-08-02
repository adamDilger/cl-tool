package entry

import (
	"bufio"
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
# Fill out the fields below (in yaml format), then save and quit

# one of (added, changed, deprecated, removed, fixed, security)
entrytype:

# the line to be added into the changelog
change:

# a short name for the file to be created (not used in final output)
filename:
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
		return fmt.Errorf("filename must not be empty")
	}

	if ce.Change == "" {
		return fmt.Errorf("change must not be empty")
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
	defer tmpFile.Close()
	tmpFile.WriteString(template)

	editor := os.ExpandEnv("$VISUAL")
	if editor == "" {
		editor = os.ExpandEnv("$EDITOR")
	}
	if editor == "" {
		editor = "vi"
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
	defer newFile.Close()

	out := make(map[string][]string)
	out[entryType] = []string{entry.Change}

	err = yaml.NewEncoder(newFile).Encode(out)
	if err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}

	fmt.Printf("%s created successfully!\n", newFile.Name())
	return nil
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
