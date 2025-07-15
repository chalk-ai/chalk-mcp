package tools

import (
	"context"
	"fmt"
	"reflect"

	"github.com/cockroachdb/errors"
	"github.com/mark3labs/mcp-go/mcp"
)

type ChalkQueryParams struct {
	ProjectRepository string            `json:"project_repository" mcp:"required,description=Path to the root of the Chalk project on disk. Should contain a chalk.yml or chalk.yaml file."`
	InputFeatures     map[string]string `json:"input_features" mcp:"description=Map of fully qualified feature names (FQNs) to their values for the query. Each key-value pair becomes --in {key}={value}."`
	OutputFeatures    []string          `json:"output_features" mcp:"description=List of fully qualified feature names (FQNs) to use as outputs for the query."`
	BranchName        string            `json:"branch_name" mcp:"description=Name of the branch to query against. If not provided, uses the mainline deployment."`
}

type ChalkQueryTool struct {
	executor CommandExecutor
}

func NewChalkQueryTool(executor CommandExecutor) *ChalkQueryTool {
	if executor == nil {
		executor = &DefaultCommandExecutor{}
	}
	return &ChalkQueryTool{
		executor: executor,
	}
}

func (t *ChalkQueryTool) Name() string {
	return "chalk_query"
}

func (t *ChalkQueryTool) Description() string {
	return "Query Chalk features using fully qualified feature names (FQNs). Supports input features, output features, and branch specification."
}

func (t *ChalkQueryTool) ParamsType() reflect.Type {
	return reflect.TypeOf(ChalkQueryParams{})
}

func (t *ChalkQueryTool) buildCommandArgs(params *ChalkQueryParams) []string {
	args := []string{"query", "--grpc"}

	for fqn, value := range params.InputFeatures {
		args = append(args, "--in", fmt.Sprintf("%s=%v", fqn, value))
	}

	for _, feature := range params.OutputFeatures {
		args = append(args, "--out", feature)
	}

	if params.BranchName != "" {
		args = append(args, "--branch="+params.BranchName)
	}

	return args
}

func (t *ChalkQueryTool) Execute(ctx context.Context, args any) (*mcp.CallToolResult, error) {
	params, ok := args.(*ChalkQueryParams)
	if !ok {
		return nil, errors.New("invalid parameter type")
	}

	output, err := t.executor.Execute(
		ctx,
		params.ProjectRepository,
		t.buildCommandArgs(params)...,
	)
	if err != nil {
		return nil, errors.Wrapf(err, "running chalk command; output: %s", output)
	}

	return mcp.NewToolResultText(string(output)), nil
}
