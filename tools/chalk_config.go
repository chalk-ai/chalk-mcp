package tools

import (
	"context"
	"reflect"

	"github.com/cockroachdb/errors"
	"github.com/mark3labs/mcp-go/mcp"
)

type ChalkConfigParams struct {
	ProjectRepository string `json:"project_repository" mcp:"required,description=Path to the root of the Chalk project on disk. Should contain a chalk.yml or chalk.yaml file."`
}

type ChalkConfigTool struct {
	executor CommandExecutor
}

func NewChalkConfigTool(executor CommandExecutor) *ChalkConfigTool {
	if executor == nil {
		executor = &DefaultCommandExecutor{}
	}
	return &ChalkConfigTool{executor: executor}
}

func (t *ChalkConfigTool) Name() string {
	return "chalk_config"
}

func (t *ChalkConfigTool) Description() string {
	return "Get the chalk config from a chalk project"
}

func (t *ChalkConfigTool) ParamsType() reflect.Type {
	return reflect.TypeOf(ChalkConfigParams{})
}

func (t *ChalkConfigTool) Execute(ctx context.Context, args any) (*mcp.CallToolResult, error) {
	params, ok := args.(*ChalkConfigParams)
	if !ok {
		return nil, errors.New("invalid parameter type")
	}

	output, err := t.executor.Execute(
		ctx,
		params.ProjectRepository,
		"config",
	)
	if err != nil {
		return nil, errors.Wrapf(err, "running chalk command; output: %s", output)
	}

	return mcp.NewToolResultText(string(output)), nil
}
