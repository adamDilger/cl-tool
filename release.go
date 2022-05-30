package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"time"

	"cl-tool/changelog"
	"cl-tool/repo"
)

func CreateRelease(root string) error {
	repo, err := repo.NewRepo(root)
	if err != nil {
		return fmt.Errorf("failed to create repository: %v", err)
	}

	isClean, err := repo.IsClean()
	if err != nil {
		return err
	}
	if !isClean {
		return fmt.Errorf("git directory has untracked changes, please stash first")
	}

	versionNumber, err := getVersionNumber()
	if err != nil {
		return err
	}

	fmt.Printf("Renaming Unreleased folder to %s - %s\n", versionNumber, time.Now().Format("2006-01-02"))
	err = renameUnreleasedFolder(root, versionNumber, time.Now().Format("2006-01-02"))
	if err != nil {
		return err
	}

	fmt.Println("Regenerating CHANGELOG.md")
	c, err := changelog.NewChangelog(root)
	if err != nil {
		return err
	}

	err = regenerateChangelogFile(root, c)
	if err != nil {
		return err
	}

	fmt.Println("Updating version in pom.xml")
	err = updateVersionInPom(root, versionNumber)
	if err != nil {
		return err
	}

	fmt.Printf("Creating branch [version-%s]\n", versionNumber)
	err = repo.CreateBranch(versionNumber)
	if err != nil {
		return err
	}

	fmt.Println("Creating commit.")
	err = repo.AddAndCommit(versionNumber)
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

	if !regexp.MustCompile(`^\d\.\d\.\d$`).MatchString(versionNumber) {
		return versionNumber, fmt.Errorf("invalid version number: %s", versionNumber)
	}

	return versionNumber, nil
}

func updateVersionInPom(path, versionNumber string) error {
	path = filepath.Join(path, "pom.xml")
	cmd := exec.Command("mvn", "--file", path, "versions:set", "-DgenerateBackupPoms=false", "-DnewVersion="+versionNumber)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to update pom.xml: %v", err)
	}
	return nil
}

func regenerateChangelogFile(path string, c *changelog.Changelog) error {
	clFile, err := os.Create(filepath.Join(path, "CHANGELOG.md"))
	if err != nil {
		return fmt.Errorf("error creating CHANGELOG.md: %v", err)
	}

	c.Render(clFile)

	return clFile.Close()
}

func renameUnreleasedFolder(path, version, date string) error {
	// clean changelog dir if any empty folders are there
	exec.Command("git", "clean", "-df", filepath.Join(path, ".changelog")).Run()

	og := filepath.Join(path, ".changelog", "Unreleased")
	new := filepath.Join(path, ".changelog", fmt.Sprintf("%s_%s", version, date))
	return os.Rename(og, new)
}
