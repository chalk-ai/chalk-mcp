package tools

import (
	"context"
	"reflect"

	"github.com/cockroachdb/errors"
	"github.com/mark3labs/mcp-go/mcp"
)

type ChalkLintParams struct {
	ProjectRepository string `json:"project_repository" mcp:"required,description=Path to the root of the Chalk project on disk to lint. Should contain a chalk.yml or chalk.yaml file."`
}

type ChalkLintTool struct {
	executor CommandExecutor
}

func NewChalkLintTool(executor CommandExecutor) *ChalkLintTool {
	if executor == nil {
		executor = &DefaultCommandExecutor{}
	}
	return &ChalkLintTool{executor: executor}
}

func (t *ChalkLintTool) Name() string {
	return "chalk_lint"
}

func (t *ChalkLintTool) Description() string {
	return "Check for errors in Chalk feature pipelines"
}

func (t *ChalkLintTool) ParamsType() reflect.Type {
	return reflect.TypeOf(ChalkLintParams{})
}

func (t *ChalkLintTool) Execute(ctx context.Context, args any) (*mcp.CallToolResult, error) {
	params, ok := args.(*ChalkLintParams)
	if !ok {
		return nil, errors.New("invalid parameter type")
	}

	cmdArgs := []string{"lint"}

	output, err := t.executor.Execute(
		ctx,
		params.ProjectRepository,
		cmdArgs...,
	)
	if err != nil {
		return nil, errors.Wrapf(err, "running chalk lint command; output: %s", output)
	}

	return mcp.NewToolResultText(string(output)), nil
}