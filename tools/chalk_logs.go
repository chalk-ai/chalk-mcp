package tools

import (
	"context"
	"reflect"

	"github.com/cockroachdb/errors"
	"github.com/mark3labs/mcp-go/mcp"
)

type ChalkLogsParams struct {
	Query             string `json:"query" mcp:"required,description=Query to search logs. Examples: 'resolver:user_features', 'component:engine message:error', 'correlation_id:abc-123'"`
	ProjectRepository string `json:"project_repository" mcp:"required,description=Path to the root of the Chalk project on disk. Should contain a chalk.yml or chalk.yaml file."`
}

type ChalkLogsTool struct {
	executor CommandExecutor
}

func NewChalkLogsTool(executor CommandExecutor) *ChalkLogsTool {
	if executor == nil {
		executor = &DefaultCommandExecutor{}
	}
	return &ChalkLogsTool{executor: executor}
}

func (t *ChalkLogsTool) Name() string {
	return "chalk_logs"
}

func (t *ChalkLogsTool) Description() string {
	return "Search through Chalk logs using powerful filtering capabilities"
}

func (t *ChalkLogsTool) ParamsType() reflect.Type {
	return reflect.TypeOf(ChalkLogsParams{})
}

func (t *ChalkLogsTool) Execute(ctx context.Context, args any) (*mcp.CallToolResult, error) {
	params, ok := args.(*ChalkLogsParams)
	if !ok {
		return nil, errors.New("invalid parameter type")
	}

	cmdArgs := []string{"logs", "--query", params.Query}

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
