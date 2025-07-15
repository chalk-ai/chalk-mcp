package tools

import (
	"reflect"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestParseParamsChalkConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		args     map[string]any
		expected *ChalkConfigParams
		wantErr  bool
	}{
		{
			name: "successful params",
			args: map[string]any{
				"project_repository": "/path/to/project",
			},
			expected: &ChalkConfigParams{
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
			expected: &ChalkConfigParams{
				ProjectRepository: "123",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			rawResult, err := parseParams(reflect.TypeOf(ChalkConfigParams{}), tt.args)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, rawResult)
				result := rawResult.(*ChalkConfigParams)
				assert.Equal(t, tt.expected.ProjectRepository, result.ProjectRepository)
			}
		})
	}
}

func TestChalkConfigTool_Execute(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		params       *ChalkConfigParams
		expectedArgs []string
		expectedOut  string
	}{
		{
			name: "successful config retrieval",
			params: &ChalkConfigParams{
				ProjectRepository: "/path/to/project",
			},
			expectedArgs: []string{"config"},
			expectedOut:  "config output",
		},
		{
			name: "different project path",
			params: &ChalkConfigParams{
				ProjectRepository: "/another/path",
			},
			expectedArgs: []string{"config"},
			expectedOut:  "different config output",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mockExecutor := new(MockCommandExecutor)
			mockExecutor.On("Execute", mock.Anything, mock.Anything, mock.Anything).Return([]byte(tt.expectedOut), nil)

			tool := NewChalkConfigTool(mockExecutor)
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
