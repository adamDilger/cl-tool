package changelog

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/go-yaml/yaml"
)

const head_default string = `# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).
`

type Changelog struct {
	Versions   []ChangelogVersion
	head, tail string
}

func NewChangelog(root string) (*Changelog, error) {
	c := Changelog{}

	path := filepath.Join(root, ".changelog")
	if err := c.readFromFiles(path); err != nil {
		return nil, err
	}

	return &c, nil
}

func (c *Changelog) readFromFiles(path string) error {
	c.head = c.parseTemplate(path, "head", head_default)
	c.tail = c.parseTemplate(path, "tail", "")

	versionFolders, err := os.ReadDir(path)
	if err != nil {
		return fmt.Errorf("failed to read changlog directory: %v", err)
	}

	for _, folder := range versionFolders {
		if !folder.IsDir() || strings.HasPrefix(folder.Name(), ".") {
			continue
		}

		v, err := c.loopFiles(path, folder.Name())
		if err != nil {
			return fmt.Errorf("failed to parse changelog files: %v", err)
		}

		if len(v.Files) == 0 {
			continue
		}

		c.Versions = append(c.Versions, v)
	}

	return nil
}

func (c *Changelog) loopFiles(path, versionFolderName string) (ChangelogVersion, error) {
	var version ChangelogVersion

	versionFiles, err := os.ReadDir(filepath.Join(path, versionFolderName))
	if err != nil {
		return version, fmt.Errorf("failed to read files from version folder %s: %v", versionFolderName, err)
	}

	// sort in reverse order
	sort.Slice(versionFiles, func(i, j int) bool { return versionFiles[i].Name() > versionFiles[j].Name() })

	r, _ := regexp.Compile(`^(\d\.\d\.\d)_(\d\d\d\d-\d\d-\d\d)`)
	if versionFolderName == "Unreleased" {
		version.Version = versionFolderName
		version.Date = ""
	} else if r.MatchString(versionFolderName) {
		matches := r.FindStringSubmatch(versionFolderName)

		version.Version = matches[1]
		version.Date = matches[2]
	} else {
		return version, fmt.Errorf("failed to parse version folder name %s: %v", versionFolderName, err)
	}

	for _, vFile := range versionFiles {
		if vFile.IsDir() || strings.HasPrefix(vFile.Name(), ".") {
			continue
		}

		path := filepath.Join(path, versionFolderName, vFile.Name())

		file, err := os.Open(path)
		if err != nil {
			return version, fmt.Errorf("failed to read %s: %v", path, err)
		}

		var clFile ChanglogFile
		err = yaml.NewDecoder(file).Decode(&clFile)
		if err != nil {
			return version, fmt.Errorf("failed to parse %s: %v", file.Name(), err)
		}

		version.Files = append(version.Files, clFile)
	}

	return version, nil
}

func (c *Changelog) parseTemplate(path, templateName, fallback string) string {
	file, err := os.Open(filepath.Join(path, templateName+".md"))
	if err != nil {
		return fallback
	}

	template, err := ioutil.ReadAll(file)
	if err != nil || len(template) == 0 {
		return fallback
	}

	return string(template)
}

func (c *Changelog) Render(writer io.Writer) {
	fmt.Fprintf(writer, "%s\n", c.head)

	//	put unreleased versions first, else sort by version number
	sort.Slice(c.Versions, func(i, j int) bool {
		if c.Versions[i].IsUnreleased() {
			return true
		} else if c.Versions[j].IsUnreleased() {
			return false
		}

		return c.Versions[i].Version > c.Versions[j].Version
	})

	for i, v := range c.Versions {
		fmt.Fprint(writer, v.Title())

		var added, changed, deprecated, removed, fixed, security EntryList

		for _, e := range v.Files {
			added = append(added, e.Added...)
			changed = append(changed, e.Changed...)
			deprecated = append(deprecated, e.Deprecated...)
			removed = append(removed, e.Removed...)
			fixed = append(fixed, e.Fixed...)
			security = append(security, e.Security...)
		}

		added.Render(writer, "Added")
		changed.Render(writer, "Changed")
		deprecated.Render(writer, "Deprecated")
		removed.Render(writer, "Removed")
		fixed.Render(writer, "Fixed")
		security.Render(writer, "Security")

		if i != len(c.Versions)-1 {
			fmt.Fprintln(writer)
		}
	}

	if c.tail != "" {
		fmt.Fprintf(writer, "\n%s", c.tail)
	}
}

type ChangelogVersion struct {
	Version, Date string
	Files         []ChanglogFile
}

func (v *ChangelogVersion) IsUnreleased() bool {
	if v.Version == "Unreleased" || v.Date == "" {
		return true
	}

	return false
}

func (v *ChangelogVersion) Title() string {
	if v.IsUnreleased() {
		return fmt.Sprintf("## [%s]\n", v.Version)
	}

	return fmt.Sprintf("## [%s] - %s\n", v.Version, v.Date)
}

type ChanglogFile struct {
	Added, Changed, Deprecated, Removed, Fixed, Security EntryList
}

type EntryList []string

func (el EntryList) Render(writer io.Writer, title string) {
	if len(el) == 0 {
		return
	}

	fmt.Fprintln(writer, "### "+title)
	for _, i := range el {
		fmt.Fprintf(writer, "- %s\n", i)
	}
}
