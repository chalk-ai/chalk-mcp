package tools

import (
	"reflect"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestParseParamsChalkLogs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		args     map[string]any
		expected *ChalkLogsParams
		wantErr  bool
	}{
		{
			name: "successful params",
			args: map[string]any{
				"query":              "resolver:user_features",
				"project_repository": "/path/to/project",
			},
			expected: &ChalkLogsParams{
				Query:             "resolver:user_features",
				ProjectRepository: "/path/to/project",
			},
		},
		{
			name: "complex query",
			args: map[string]any{
				"query":              "component:engine message:error",
				"project_repository": "/path/to/project",
			},
			expected: &ChalkLogsParams{
				Query:             "component:engine message:error",
				ProjectRepository: "/path/to/project",
			},
		},
		{
			name: "correlation_id query",
			args: map[string]any{
				"query":              "correlation_id:abc-123",
				"project_repository": "/path/to/project",
			},
			expected: &ChalkLogsParams{
				Query:             "correlation_id:abc-123",
				ProjectRepository: "/path/to/project",
			},
		},
		{
			name: "missing query",
			args: map[string]any{
				"project_repository": "/path/to/project",
			},
			wantErr: true,
		},
		{
			name: "missing project repository",
			args: map[string]any{
				"query": "resolver:user_features",
			},
			wantErr: true,
		},
		{
			name:    "missing both parameters",
			args:    map[string]any{},
			wantErr: true,
		},
		{
			name: "coalesced type for query",
			args: map[string]any{
				"query":              123,
				"project_repository": "/path/to/project",
			},
			expected: &ChalkLogsParams{
				Query:             "123",
				ProjectRepository: "/path/to/project",
			},
		},
		{
			name: "coalesced type for project repository",
			args: map[string]any{
				"query":              "resolver:user_features",
				"project_repository": 456,
			},
			expected: &ChalkLogsParams{
				Query:             "resolver:user_features",
				ProjectRepository: "456",
			},
		},
		{
			name: "query with special characters",
			args: map[string]any{
				"query":              "message:\"error occurred\" component:engine",
				"project_repository": "/path/to/project",
			},
			expected: &ChalkLogsParams{
				Query:             "message:\"error occurred\" component:engine",
				ProjectRepository: "/path/to/project",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			rawResult, err := parseParams(reflect.TypeOf(ChalkLogsParams{}), tt.args)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, rawResult)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, rawResult)
				result := rawResult.(*ChalkLogsParams)
				assert.Equal(t, tt.expected.Query, result.Query)
				assert.Equal(t, tt.expected.ProjectRepository, result.ProjectRepository)
			}
		})
	}
}

func TestChalkLogsTool_Execute(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		params       *ChalkLogsParams
		expectedArgs []string
		expectedOut  string
	}{
		{
			name: "simple resolver query",
			params: &ChalkLogsParams{
				Query:             "resolver:user_features",
				ProjectRepository: "/path/to/project",
			},
			expectedArgs: []string{"logs", "--query", "resolver:user_features"},
			expectedOut:  "2024-01-01 10:00:00 [INFO] resolver:user_features - Feature computed successfully",
		},
		{
			name: "component and message query",
			params: &ChalkLogsParams{
				Query:             "component:engine message:error",
				ProjectRepository: "/path/to/project",
			},
			expectedArgs: []string{"logs", "--query", "component:engine message:error"},
			expectedOut:  "2024-01-01 10:01:00 [ERROR] component:engine - Error processing request",
		},
		{
			name: "correlation_id query",
			params: &ChalkLogsParams{
				Query:             "correlation_id:abc-123",
				ProjectRepository: "/another/path",
			},
			expectedArgs: []string{"logs", "--query", "correlation_id:abc-123"},
			expectedOut:  "2024-01-01 10:02:00 [INFO] correlation_id:abc-123 - Request processed",
		},
		{
			name: "query with special characters",
			params: &ChalkLogsParams{
				Query:             "message:\"error occurred\" component:engine",
				ProjectRepository: "/path/to/project",
			},
			expectedArgs: []string{"logs", "--query", "message:\"error occurred\" component:engine"},
			expectedOut:  "2024-01-01 10:03:00 [ERROR] component:engine - Error occurred in processing",
		},
		{
			name: "empty query",
			params: &ChalkLogsParams{
				Query:             "",
				ProjectRepository: "/path/to/project",
			},
			expectedArgs: []string{"logs", "--query", ""},
			expectedOut:  "2024-01-01 10:04:00 [INFO] - All logs",
		},
		{
			name: "project with spaces",
			params: &ChalkLogsParams{
				Query:             "resolver:user_features",
				ProjectRepository: "/path/with spaces/project",
			},
			expectedArgs: []string{"logs", "--query", "resolver:user_features"},
			expectedOut:  "2024-01-01 10:05:00 [INFO] resolver:user_features - Feature computed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mockExecutor := new(MockCommandExecutor)
			mockExecutor.On("Execute", mock.Anything, mock.Anything, mock.Anything).Return([]byte(tt.expectedOut), nil)

			tool := NewChalkLogsTool(mockExecutor)
			result, err := tool.Execute(t.Context(), tt.params)
			assert.NoError(t, err)
			assert.Equal(t, mcp.NewToolResultText(tt.expectedOut), result)

			mockExecutor.AssertCalled(t, "Execute", mock.Anything, tt.params.ProjectRepository, mock.Anything)
			assert.Len(t, mockExecutor.Calls, 1)
			actualArgs := mockExecutor.Calls[0].Arguments[2].([]string)
			assert.Equal(t, tt.expectedArgs, actualArgs)
		})
	}
}
