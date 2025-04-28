package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

func TestChalkConfigHandler(t *testing.T) {
	tests := []struct {
		name              string
		projectRepository string
		setup             func()
		expectedError     error
	}{
		{
			name:              "Valid project repository",
			projectRepository: "/Users/andrew/chalk/fraud-template-staging",
			setup: func() {
				// Ensure the directory and chalk.yml exist
				//os.MkdirAll("/Users/andrew/fraud-template-staging", 0755)
				//os.WriteFile("/Users/andrew/fraud-template-staging/chalk.yml", []byte("config: test"), 0644)
			},
			expectedError: nil,
		},
		{
			name:              "Missing project repository",
			projectRepository: "/Users/andrew/nonexistent",
			setup:             func() {}, // No setup, directory doesn't exist
			expectedError:     errors.New("project_repository must exist"),
		},
		{
			name:              "Missing chalk.yml",
			projectRepository: "/Users/andrew/fraud-template-staging",
			setup: func() {
				// Ensure the directory exists but no chalk.yml
				os.MkdirAll("/Users/andrew/fraud-template-staging", 0755)
				os.Remove("/Users/andrew/fraud-template-staging/chalk.yml")
			},
			expectedError: errors.New("project_repository must contain a chalk.yml file"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

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
