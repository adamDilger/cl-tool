package repo

import (
	"os/exec"
)

type Repo struct {
	path string
}

func NewRepo(path string) (*Repo, error) {
	return &Repo{path: path}, nil
}

func (r *Repo) IsClean() (bool, error) {
	cmd := exec.Command("git", "-C", r.path, "status", "--porcelain")

	o, err := cmd.Output()
	if err != nil {
		return false, err
	}

	return string(o) == "", nil
}

func (r *Repo) CreateBranch(versionNumber string) error {
	branchName := "version-" + versionNumber
	cmd := exec.Command("git", "-C", r.path, "checkout", "-b", branchName)

	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func (r *Repo) AddAndCommit(versionNumber string) error {
	cmd := exec.Command("git", "-C", r.path, "add", "pom.xml", "CHANGELOG.md", ".changelog")

	err := cmd.Run()
	if err != nil {
		return err
	}

	cmd = exec.Command("git", "-C", r.path, "commit", "-m", "Version "+versionNumber)

	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
