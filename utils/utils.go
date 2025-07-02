package utils

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
)

// HasChalkConfig checks if a directory contains either chalk.yml or chalk.yaml
func HasChalkConfig(dir string) (bool, error) {
	if _, err := os.Stat(filepath.Join(dir, "chalk.yml")); err == nil {
		return true, nil
	}
	if _, err := os.Stat(filepath.Join(dir, "chalk.yaml")); err == nil {
		return true, nil
	}
	return false, nil
}

// FindChalkBinary attempts to locate the chalk binary in common locations
func FindChalkBinary() (string, error) {
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

	for _, path := range commonPaths {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	return "", errors.New("chalk binary not found in PATH or common locations")
}
