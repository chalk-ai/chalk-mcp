package utils

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
)

// hasChalkConfig checks if a directory contains either chalk.yml or chalk.yaml
func hasChalkConfig(dir string) bool {
	if _, err := os.Stat(filepath.Join(dir, "chalk.yml")); err == nil {
		return true
	}
	if _, err := os.Stat(filepath.Join(dir, "chalk.yaml")); err == nil {
		return true
	}
	return false
}

// validateChalkProject checks if a repository contains a chalk.yml or chalk.yaml file
func validateChalkProject(projectRepository string) error {
	if _, err := os.Stat(projectRepository); os.IsNotExist(err) {
		return errors.New("project_repository must exist")
	}

	hasConfig := hasChalkConfig(projectRepository)
	if !hasConfig {
		return errors.New("project_repository must contain a chalk.yml or chalk.yaml file")
	}

	return nil
}

// findChalkBinary attempts to locate the chalk binary in common locations
func findChalkBinary(checkCommonPaths bool) (string, error) {
	if path, err := exec.LookPath("chalk"); err == nil {
		return path, nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	commonPaths := []string{
		filepath.Join(homeDir, ".chalk", "bin", "chalk-latest"),
		filepath.Join(homeDir, ".chalk", "bin", "chalk"),
		"/usr/local/bin/chalk",
		"/opt/chalk/bin/chalk",
	}

	if checkCommonPaths {
		for _, path := range commonPaths {
			if _, err := os.Stat(path); err == nil {
				return path, nil
			}
		}
	}
	return "", errors.New("chalk binary not found in PATH or common locations")
}

func GetChalkCommand(projectRepository string, args ...string) (*exec.Cmd, error) {
	if err := validateChalkProject(projectRepository); err != nil {
		return nil, err
	}

	chalkPath, err := findChalkBinary(true)
	if err != nil {
		return nil, err
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	cmd := exec.Command(chalkPath, args...)
	cmd.Dir = projectRepository
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "XDG_CONFIG_HOME="+homeDir+"/")
	return cmd, nil
}
