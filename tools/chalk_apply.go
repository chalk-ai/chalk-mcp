package tools

import (
	"context"
	"reflect"

	"github.com/cockroachdb/errors"
	"github.com/mark3labs/mcp-go/mcp"
)

type ChalkApplyParams struct {
	ProjectRepository string `json:"project_repository" mcp:"required,description=Path to the root of the Chalk project on disk. Should contain a chalk.yml or chalk.yaml file."`
	BranchName        string `json:"branch_name" mcp:"description=Name of the branch to deploy to. If not provided, defaults to the current git branch."`
}

type ChalkApplyTool struct {
	executor CommandExecutor
}

func NewChalkApplyTool(executor CommandExecutor) *ChalkApplyTool {
	if executor == nil {
		executor = &DefaultCommandExecutor{}
	}
	return &ChalkApplyTool{executor: executor}
}

func (t *ChalkApplyTool) Name() string {
	return "chalk_apply"
}

func (t *ChalkApplyTool) Description() string {
	return "Apply Chalk configurations with branch and JSON output. Recommended workflow is to use the chalk_environment tool to set the environment, then use this tool to apply the configuration."
}

func (t *ChalkApplyTool) ParamsType() reflect.Type {
	return reflect.TypeOf(ChalkApplyParams{})
}

func (t *ChalkApplyTool) Execute(ctx context.Context, args any) (*mcp.CallToolResult, error) {
	params, ok := args.(*ChalkApplyParams)
	if !ok {
		return nil, errors.New("invalid parameter type")
	}

	cmdArgs := []string{"apply"}
	if params.BranchName != "" {
		cmdArgs = append(cmdArgs, "--branch="+params.BranchName)
	} else {
		cmdArgs = append(cmdArgs, "--branch")
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
