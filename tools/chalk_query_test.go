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
	tests := []struct {
		name          string
		params        *ChalkQueryParams
		expectedArgs  []string
		expectedOut   string
		expectedError bool
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
				"--in", "user.id=456",
				"--out", "user.email",
			},
			expectedOut: "query result without branch",
		},
		{
			name: "test arg validation",
			params: &ChalkQueryParams{
				ProjectRepository: "/path/to/project",
				InputFeatures: map[string]string{
					"user.id": "456",
				},
				OutputFeatures: []string{"user.email"},
			},
			expectedArgs:  []string{"--in", "user.id=456", "--out", "user.email"}, // Expected args to validate
			expectedOut:   "query result for validation test",
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockExecutor := new(MockCommandExecutor)
			mockExecutor.On("Execute", mock.Anything, mock.Anything, mock.Anything).Return([]byte(tt.expectedOut), nil)

			tool := NewChalkQueryTool(mockExecutor)
			result, err := tool.Execute(t.Context(), tt.params)
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, mcp.NewToolResultText(tt.expectedOut), result)

				mockExecutor.AssertCalled(t, "Execute", mock.Anything, tt.params.ProjectRepository, mock.Anything)
				assert.Len(t, mockExecutor.Calls, 1)
				actualArgs := mockExecutor.Calls[0].Arguments[2].([]string)
				assert.Equal(t, append([]string{"query", "--grpc"}, tt.expectedArgs...), actualArgs)
			}
		})
	}
}

func TestGenerateMetadataForChalkQueryTool(t *testing.T) {
	tool := NewChalkQueryTool(nil)
	metadata := GenerateMetadata(tool)

	expected := mcp.NewTool("chalk_query",
		mcp.WithDescription("Query Chalk features using fully qualified feature names (FQNs). Supports input features, output features, and branch specification."),
		mcp.WithString("project_repository",
			mcp.Required(),
			mcp.Description("Path to the root of the Chalk project on disk. Should contain a chalk.yml or chalk.yaml file."),
		),
		mcp.WithObject("input_features",
			mcp.Description("Map of fully qualified feature names (FQNs) to their values for the query. Each key-value pair becomes --in {key}={value}."),
			mcp.AdditionalProperties(map[string]any{"type": "string"}),
		),
		mcp.WithArray("output_features",
			mcp.Description("List of fully qualified feature names (FQNs) to use as outputs for the query."),
			mcp.Items(map[string]any{"type": "string"}),
		),
		mcp.WithString("branch_name",
			mcp.Description("Name of the branch to query against. If not provided, uses the mainline deployment."),
		),
	)
	assert.Equal(t, expected, metadata)
}
