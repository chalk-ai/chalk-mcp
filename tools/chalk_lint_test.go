package tools

import (
	"reflect"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestParseParamsChalkLint(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		args     map[string]any
		expected *ChalkLintParams
		wantErr  bool
	}{
		{
			name: "successful params",
			args: map[string]any{
				"project_repository": "/path/to/project",
			},
			expected: &ChalkLintParams{
				ProjectRepository: "/path/to/project",
			},
		},
		{
			name:    "missing project repository",
			args:    map[string]any{},
			wantErr: true,
		},
		{
			name: "coalesced type for project repository",
			args: map[string]any{
				"project_repository": 123,
			},
			expected: &ChalkLintParams{
				ProjectRepository: "123",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			rawResult, err := parseParams(reflect.TypeOf(ChalkLintParams{}), tt.args)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, rawResult)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, rawResult)
				result := rawResult.(*ChalkLintParams)
				assert.Equal(t, tt.expected.ProjectRepository, result.ProjectRepository)
			}
		})
	}
}

func TestChalkLintTool_Execute(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		params       *ChalkLintParams
		expectedArgs []string
		expectedOut  string
	}{
		{
			name: "successful lint check",
			params: &ChalkLintParams{
				ProjectRepository: "/path/to/project",
			},
			expectedArgs: []string{"lint"},
			expectedOut:  "No errors found in feature pipelines",
		},
		{
			name: "lint with errors",
			params: &ChalkLintParams{
				ProjectRepository: "/another/path",
			},
			expectedArgs: []string{"lint"},
			expectedOut:  "Error: feature validation failed\nError: type mismatch in feature definition",
		},
		{
			name: "project with special characters",
			params: &ChalkLintParams{
				ProjectRepository: "/path/with spaces/project",
			},
			expectedArgs: []string{"lint"},
			expectedOut:  "Warning: deprecated feature usage found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mockExecutor := new(MockCommandExecutor)
			mockExecutor.On("Execute", mock.Anything, mock.Anything, mock.Anything).Return([]byte(tt.expectedOut), nil)

			tool := NewChalkLintTool(mockExecutor)
			result, err := tool.Execute(t.Context(), tt.params)
			assert.NoError(t, err)
			assert.Equal(t, mcp.NewToolResultText(tt.expectedOut), result)

			mockExecutor.AssertCalled(t, "Execute", mock.Anything, tt.params.ProjectRepository, tt.expectedArgs)
			assert.Len(t, mockExecutor.Calls, 1)
		})
	}
}
