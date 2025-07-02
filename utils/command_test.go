package utils

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createTestTempDir(t *testing.T) string {
	tempDir, err := os.MkdirTemp("", "chalk-test-*")
	assert.NoError(t, err)
	assert.DirExists(t, tempDir)
	return tempDir
}

func createTestFile(t *testing.T, tempDir string, fileName string, content string) {
	filePath := filepath.Join(tempDir, fileName)
	err := os.WriteFile(filePath, []byte(content), 0644)
	assert.NoError(t, err)
	assert.FileExists(t, filePath)
}

func createTestChalkBinary(t *testing.T, tempDir string) {
	chalkBinary := filepath.Join(tempDir, "chalk")
	err := os.WriteFile(
		chalkBinary,
		[]byte("#!/bin/sh\necho \"$@\"\n"),
		0755,
	)
	assert.NoError(t, err)
	assert.FileExists(t, chalkBinary)
}

func TestValidateChalkProject(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		createTmpDir bool
		tmpFiles     []string
		expectedErr  error
	}{
		{
			name:         "valid project repository with chalk.yml",
			createTmpDir: true,
			tmpFiles:     []string{"chalk.yml"},
			expectedErr:  nil,
		},
		{
			name:         "valid project repository with chalk.yaml",
			createTmpDir: true,
			tmpFiles:     []string{"chalk.yaml"},
			expectedErr:  nil,
		},
		{
			name:         "missing project repository",
			createTmpDir: false,
			tmpFiles:     []string{},
			expectedErr:  errors.New("project_repository must exist"),
		},
		{
			name:         "missing chalk yaml",
			createTmpDir: true,
			tmpFiles:     []string{"not-chalk.yaml"},
			expectedErr:  errors.New("project_repository must contain a chalk.yml or chalk.yaml file"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var tempDir string
			if tt.createTmpDir {
				tempDir = createTestTempDir(t)
			}
			defer func() {
				if tt.createTmpDir {
					os.RemoveAll(tempDir)
				}
			}()

			for _, tmpFile := range tt.tmpFiles {
				createTestFile(t, tempDir, tmpFile, "config: test")
			}

			err := validateChalkProject(tempDir)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestFindChalkBinary(t *testing.T) {
	t.Parallel()

	t.Run("in PATH", func(t *testing.T) {
		t.Parallel()

		tempDir := createTestTempDir(t)
		defer os.RemoveAll(tempDir)

		createTestChalkBinary(t, tempDir)
		os.Setenv("PATH", tempDir)

		chalkPath, err := findChalkBinary(false)
		assert.NoError(t, err)
		assert.NotEmpty(t, chalkPath)
	})

	t.Run("not in common paths", func(t *testing.T) {
		t.Parallel()

		os.Setenv("PATH", "")

		chalkPath, err := findChalkBinary(false)
		assert.Error(t, err)
		assert.Empty(t, chalkPath)
	})
}

func TestGetChalkCommand(t *testing.T) {
	t.Parallel()

	tempDir := createTestTempDir(t)
	defer os.RemoveAll(tempDir)

	createTestFile(t, tempDir, "chalk.yml", "config: test")

	createTestChalkBinary(t, tempDir)
	os.Setenv("PATH", tempDir)

	projectRepository := tempDir
	testArg := "my_favorite_arg"
	cmd, err := GetChalkCommand(projectRepository, testArg)
	assert.NoError(t, err)

	out, err := cmd.CombinedOutput()
	assert.NoError(t, err)
	assert.Equal(t, testArg+"\n", string(out))
}
