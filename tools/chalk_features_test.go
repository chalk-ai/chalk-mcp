package tools

import (
	"reflect"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestParseParamsChalkFeatures(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		args     map[string]any
		expected *ChalkFeaturesParams
		wantErr  bool
	}{
		{
			name: "successful params",
			args: map[string]any{
				"project_repository": "/path/to/project",
			},
			expected: &ChalkFeaturesParams{
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
			expected: &ChalkFeaturesParams{
				ProjectRepository: "123",
			},
		},
		{
			name: "project repository with special characters",
			args: map[string]any{
				"project_repository": "/path/to/project with spaces/",
			},
			expected: &ChalkFeaturesParams{
				ProjectRepository: "/path/to/project with spaces/",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			rawResult, err := parseParams(reflect.TypeOf(ChalkFeaturesParams{}), tt.args)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, rawResult)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, rawResult)
				result := rawResult.(*ChalkFeaturesParams)
				assert.Equal(t, tt.expected.ProjectRepository, result.ProjectRepository)
			}
		})
	}
}

func TestChalkFeaturesTool_Execute(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		params       *ChalkFeaturesParams
		expectedArgs []string
		expectedOut  string
	}{
		{
			name: "successful features retrieval",
			params: &ChalkFeaturesParams{
				ProjectRepository: "/path/to/project",
			},
			expectedArgs: []string{"features"},
			expectedOut:  "feature1\nfeature2\nfeature3",
		},
		{
			name: "different project path",
			params: &ChalkFeaturesParams{
				ProjectRepository: "/another/path",
			},
			expectedArgs: []string{"features"},
			expectedOut:  "user.id\nuser.name\nuser.email",
		},
		{
			name: "project with special characters",
			params: &ChalkFeaturesParams{
				ProjectRepository: "/path/with spaces/project",
			},
			expectedArgs: []string{"features"},
			expectedOut:  "special.feature1\nspecial.feature2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mockExecutor := new(MockCommandExecutor)
			mockExecutor.On("Execute", mock.Anything, mock.Anything, mock.Anything).Return([]byte(tt.expectedOut), nil)

			tool := NewChalkFeaturesTool(mockExecutor)
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
