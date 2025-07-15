package tools

import (
	"context"
	"reflect"

	"github.com/cockroachdb/errors"
	"github.com/mark3labs/mcp-go/mcp"
)

type ChalkEnvironmentOperation string

const (
	ChalkEnvironmentOperationSet ChalkEnvironmentOperation = "set"
	ChalkEnvironmentOperationGet ChalkEnvironmentOperation = "get"
)

type ChalkEnvironmentParams struct {
	ProjectRepository string `json:"project_repository" mcp:"required,description=Path to the root of the Chalk project on disk. Should contain a chalk.yml or chalk.yaml file."`
	Operation         string `json:"operation" mcp:"required,description=The operation to perform on the chalk environment"`
	Environment       string `json:"environment" mcp:"description=If doing a set operation, the environment to set for the current project. Can be an environment ID or name."`
}

type ChalkEnvironmentTool struct {
	executor CommandExecutor
}

func NewChalkEnvironmentTool(executor CommandExecutor) *ChalkEnvironmentTool {
	if executor == nil {
		executor = &DefaultCommandExecutor{}
	}
	return &ChalkEnvironmentTool{executor: executor}
}

func (t *ChalkEnvironmentTool) Name() string {
	return "chalk_environment"
}

func (t *ChalkEnvironmentTool) Description() string {
	return "Get or set the chalk environment for the current project"
}

func (t *ChalkEnvironmentTool) ParamsType() reflect.Type {
	return reflect.TypeOf(ChalkEnvironmentParams{})
}

func (t *ChalkEnvironmentTool) Execute(ctx context.Context, args any) (*mcp.CallToolResult, error) {
	params, ok := args.(*ChalkEnvironmentParams)
	if !ok {
		return nil, errors.New("invalid parameter type")
	}

	var cmdArgs []string
	switch params.Operation {
	case string(ChalkEnvironmentOperationSet):
		if params.Environment == "" {
			return nil, errors.New("environment must be provided for set operation")
		}
		cmdArgs = []string{"environment", params.Environment}
	case string(ChalkEnvironmentOperationGet):
		cmdArgs = []string{"environment"}
	default:
		return nil, errors.Newf("invalid operation: %s", params.Operation)
	}

	output, err := t.executor.Execute(
		ctx,
		params.ProjectRepository,
		cmdArgs...,
	)
	if err != nil {
		return nil, errors.Wrapf(err, "running chalk command; output: %s", output)
	}

	return mcp.NewToolResultText(string(output)), nil
}
