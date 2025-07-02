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
	// Create a temporary directory for valid test case
	validTmpDir, err := os.MkdirTemp("", "chalk-test-valid-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(validTmpDir)

	// Create chalk.yml in the valid temp directory
	chalkYml := filepath.Join(validTmpDir, "chalk.yml")
	if err := os.WriteFile(chalkYml, []byte("config: test"), 0644); err != nil {
		t.Fatalf("Failed to write chalk.yml: %v", err)
	}

	// Create a temporary directory for missing chalk.yml test case
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
			name:              "Valid project repository",
			projectRepository: validTmpDir,
			expectedError:     nil,
		},
		{
			name:              "Missing project repository",
			projectRepository: "/tmp/nonexistent-chalk-test",
			expectedError:     errors.New("project_repository must exist"),
		},
		{
			name:              "Missing chalk.yml",
			projectRepository: invalidTmpDir,
			expectedError:     errors.New("project_repository must contain a chalk.yml file"),
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