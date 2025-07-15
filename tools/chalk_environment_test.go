package tools

import (
	"reflect"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestParseParamsChalkEnvironment(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		args     map[string]any
		expected *ChalkEnvironmentParams
		wantErr  bool
	}{
		{
			name: "successful set operation",
			args: map[string]any{
				"project_repository": "/path/to/project",
				"operation":          "set",
				"environment":        "production",
			},
			expected: &ChalkEnvironmentParams{
				ProjectRepository: "/path/to/project",
				Operation:         "set",
				Environment:       "production",
			},
		},
		{
			name: "successful get operation",
			args: map[string]any{
				"project_repository": "/path/to/project",
				"operation":          "get",
			},
			expected: &ChalkEnvironmentParams{
				ProjectRepository: "/path/to/project",
				Operation:         "get",
				Environment:       "",
			},
		},
		{
			name: "missing operation",
			args: map[string]any{
				"project_repository": "/path/to/project",
			},
			wantErr: true,
		},
		{
			name: "missing project repository",
			args: map[string]any{
				"operation": "get",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			rawResult, err := parseParams(reflect.TypeOf(ChalkEnvironmentParams{}), tt.args)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, rawResult)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, rawResult)
				result := rawResult.(*ChalkEnvironmentParams)
				assert.Equal(t, tt.expected.ProjectRepository, result.ProjectRepository)
				assert.Equal(t, tt.expected.Operation, result.Operation)
				assert.Equal(t, tt.expected.Environment, result.Environment)
			}
		})
	}
}

func TestChalkEnvironmentTool_Execute(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		params       *ChalkEnvironmentParams
		expectedArgs []string
		expectedOut  string
		expectedErr  bool
	}{
		{
			name: "set operation",
			params: &ChalkEnvironmentParams{
				ProjectRepository: "/path/to/project",
				Operation:         "set",
				Environment:       "production",
			},
			expectedArgs: []string{"environment", "production"},
			expectedOut:  "Environment set to production",
		},
		{
			name: "get operation",
			params: &ChalkEnvironmentParams{
				ProjectRepository: "/path/to/project",
				Operation:         "get",
			},
			expectedArgs: []string{"environment"},
			expectedOut:  "Current environment: staging",
		},
		{
			name: "invalid operation",
			params: &ChalkEnvironmentParams{
				ProjectRepository: "/path/to/project",
				Operation:         "invalid",
			},
			expectedErr: true,
		},
		{
			name: "set operation without environment",
			params: &ChalkEnvironmentParams{
				ProjectRepository: "/path/to/project",
				Operation:         "set",
				Environment:       "",
			},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mockExecutor := new(MockCommandExecutor)

			if !tt.expectedErr {
				mockExecutor.On("Execute", mock.Anything, mock.Anything, mock.Anything).Return([]byte(tt.expectedOut), nil)
			}

			tool := NewChalkEnvironmentTool(mockExecutor)
			result, err := tool.Execute(t.Context(), tt.params)

			if tt.expectedErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, mcp.NewToolResultText(tt.expectedOut), result)

				mockExecutor.AssertCalled(t, "Execute", mock.Anything, tt.params.ProjectRepository, mock.Anything)
				assert.Len(t, mockExecutor.Calls, 1)
				actualArgs := mockExecutor.Calls[0].Arguments[2].([]string)
				assert.Equal(t, tt.expectedArgs, actualArgs)
			}
		})
	}
}
