package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

func TestChalkConfigHandler(t *testing.T) {
	// Create a temporary directory for valid test case with chalk.yml
	validTmpDirYml, err := os.MkdirTemp("", "chalk-test-valid-yml-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(validTmpDirYml)

	// Create chalk.yml in the valid temp directory
	chalkYml := filepath.Join(validTmpDirYml, "chalk.yml")
	if err := os.WriteFile(chalkYml, []byte("config: test"), 0644); err != nil {
		t.Fatalf("Failed to write chalk.yml: %v", err)
	}

	// Create a temporary directory for valid test case with chalk.yaml
	validTmpDirYaml, err := os.MkdirTemp("", "chalk-test-valid-yaml-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(validTmpDirYaml)

	// Create chalk.yaml in the valid temp directory
	chalkYaml := filepath.Join(validTmpDirYaml, "chalk.yaml")
	if err := os.WriteFile(chalkYaml, []byte("config: test"), 0644); err != nil {
		t.Fatalf("Failed to write chalk.yaml: %v", err)
	}

	// Create a temporary directory for missing chalk config test case
	invalidTmpDir, err := os.MkdirTemp("", "chalk-test-invalid-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(invalidTmpDir)

	tests := []struct {
		name              string
		projectRepository string
		expectedError     error
	}{
		{
			name:              "Valid project repository with chalk.yml",
			projectRepository: validTmpDirYml,
			expectedError:     nil,
		},
		{
			name:              "Valid project repository with chalk.yaml",
			projectRepository: validTmpDirYaml,
			expectedError:     nil,
		},
		{
			name:              "Missing project repository",
			projectRepository: "/tmp/nonexistent-chalk-test",
			expectedError:     errors.New("project_repository must exist"),
		},
		{
			name:              "Missing chalk config file",
			projectRepository: invalidTmpDir,
			expectedError:     errors.New("project_repository must contain a chalk.yml or chalk.yaml file"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := mcp.CallToolRequest{}
			request.Params.Arguments = make(map[string]interface{})
			request.Params.Arguments["project_repository"] = tt.projectRepository

			c, err := configHandler(context.Background(), request)
			if (err != nil && tt.expectedError == nil) || (err == nil && tt.expectedError != nil) || (err != nil && err.Error() != tt.expectedError.Error()) {
				t.Errorf("expected error: %v, got: %v", tt.expectedError, err)
			}

			fmt.Println("Config: ", c)
		})
	}
}