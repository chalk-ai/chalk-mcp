package tools

import (
	"reflect"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestParseParamsChalkQuery(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		args     map[string]any
		expected *ChalkQueryParams
		wantErr  bool
	}{
		{
			name: "successful params no branch",
			args: map[string]any{
				"project_repository": "/path/to/project",
				"input_features": map[string]any{
					"user.id":  "123",
					"user.age": 25,
				},
				"output_features": []any{"user.name", "user.email"},
			},
			expected: &ChalkQueryParams{
				ProjectRepository: "/path/to/project",
				InputFeatures: map[string]string{
					"user.id":  "123",
					"user.age": "25",
				},
				OutputFeatures: []string{"user.name", "user.email"},
			},
		},
		{
			name: "successful params with branch",
			args: map[string]any{
				"project_repository": "/path/to/project",
				"input_features": map[string]any{
					"user.id":  "123",
					"user.age": 25,
				},
				"output_features": []any{"user.name", "user.email"},
				"branch_name":     "main",
			},
			expected: &ChalkQueryParams{
				ProjectRepository: "/path/to/project",
				InputFeatures: map[string]string{
					"user.id":  "123",
					"user.age": "25",
				},
				OutputFeatures: []string{"user.name", "user.email"},
				BranchName:     "main",
			},
		},
		{
			name: "missing project repository",
			args: map[string]any{
				"input_features": map[string]any{
					"user.id":  "123",
					"user.age": 25,
				},
				"output_features": []any{"user.name", "user.email"},
			},
			wantErr: true,
		},
		{
			name: "missing output features",
			args: map[string]any{
				"project_repository": "/path/to/project",
				"input_features": map[string]any{
					"user.id":  "123",
					"user.age": 25,
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			rawResult, err := parseParams(reflect.TypeOf(ChalkQueryParams{}), tt.args)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, rawResult)
				result := rawResult.(*ChalkQueryParams)
				assert.Equal(t, tt.expected.ProjectRepository, result.ProjectRepository)
				assert.Subset(t, tt.expected.InputFeatures, result.InputFeatures)
				assert.Subset(t, result.InputFeatures, tt.expected.InputFeatures)
				assert.ElementsMatch(t, tt.expected.OutputFeatures, result.OutputFeatures)
				assert.Equal(t, tt.expected.BranchName, result.BranchName)
			}
		})
	}
}

func TestChalkQueryTool_Execute(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		params       *ChalkQueryParams
		expectedArgs []string
		expectedOut  string
	}{
		{
			name: "successful query with all parameters",
			params: &ChalkQueryParams{
				ProjectRepository: "/path/to/project",
				InputFeatures: map[string]string{
					"user.id":  "123",
					"user.age": "25",
				},
				OutputFeatures: []string{"user.name", "user.email"},
				BranchName:     "feature-branch",
			},
			expectedArgs: []string{
				"query",
				"--grpc",
				"--in", "user.id=123",
				"--in", "user.age=25",
				"--out", "user.name",
				"--out", "user.email",
				"--branch=feature-branch",
			},
			expectedOut: "successful query result",
		},
		{
			name: "query without input features",
			params: &ChalkQueryParams{
				ProjectRepository: "/path/to/project",
				OutputFeatures:    []string{"user.name"},
				BranchName:        "main",
			},
			expectedArgs: []string{
				"query",
				"--grpc",
				"--out", "user.name",
				"--branch=main",
			},
			expectedOut: "query result without inputs",
		},
		{
			name: "query without branch name",
			params: &ChalkQueryParams{
				ProjectRepository: "/path/to/project",
				InputFeatures: map[string]string{
					"user.id": "456",
				},
				OutputFeatures: []string{"user.email"},
			},
			expectedArgs: []string{
				"query",
				"--grpc",
				"--in", "user.id=456",
				"--out", "user.email",
			},
			expectedOut: "query result without branch",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mockExecutor := new(MockCommandExecutor)
			mockExecutor.On("Execute", mock.Anything, mock.Anything, mock.Anything).Return([]byte(tt.expectedOut), nil)

			tool := NewChalkQueryTool(mockExecutor)
			result, err := tool.Execute(t.Context(), tt.params)
			assert.NoError(t, err)
			assert.Equal(t, mcp.NewToolResultText(tt.expectedOut), result)

			mockExecutor.AssertCalled(t, "Execute", mock.Anything, tt.params.ProjectRepository, mock.Anything)
			assert.Len(t, mockExecutor.Calls, 1)
			actualArgs := mockExecutor.Calls[0].Arguments[2].([]string)
			// elements match since order is not guaranteed for input features
			assert.ElementsMatch(t, tt.expectedArgs, actualArgs)
		})
	}
}
