package tools

import (
	"reflect"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestParseParamsChalkApply(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		args     map[string]any
		expected *ChalkApplyParams
		wantErr  bool
	}{
		{
			name: "successful params without branch",
			args: map[string]any{
				"project_repository": "/path/to/project",
			},
			expected: &ChalkApplyParams{
				ProjectRepository: "/path/to/project",
			},
		},
		{
			name: "successful params with branch",
			args: map[string]any{
				"project_repository": "/path/to/project",
				"branch_name":        "main",
			},
			expected: &ChalkApplyParams{
				ProjectRepository: "/path/to/project",
				BranchName:        "main",
			},
		},
		{
			name:    "missing project repository",
			args:    map[string]any{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			rawResult, err := parseParams(reflect.TypeOf(ChalkApplyParams{}), tt.args)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, rawResult)
				result := rawResult.(*ChalkApplyParams)
				assert.Equal(t, tt.expected.ProjectRepository, result.ProjectRepository)
				assert.Equal(t, tt.expected.BranchName, result.BranchName)
			}
		})
	}
}

func TestChalkApplyTool_Execute(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		params       *ChalkApplyParams
		expectedArgs []string
		expectedOut  string
	}{
		{
			name: "successful apply with branch",
			params: &ChalkApplyParams{
				ProjectRepository: "/path/to/project",
				BranchName:        "main",
			},
			expectedArgs: []string{"apply", "--branch=main"},
			expectedOut:  "apply result with branch",
		},
		{
			name: "successful apply without branch",
			params: &ChalkApplyParams{
				ProjectRepository: "/path/to/project",
			},
			expectedArgs: []string{"apply", "--branch"},
			expectedOut:  "apply result without branch",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mockExecutor := new(MockCommandExecutor)
			mockExecutor.On("Execute", mock.Anything, mock.Anything, mock.Anything).Return([]byte(tt.expectedOut), nil)

			tool := NewChalkApplyTool(mockExecutor)
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