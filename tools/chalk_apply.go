package tools

import (
	"context"

	"github.com/chalk-ai/chalk-mcp/utils"
	"github.com/cockroachdb/errors"
	"github.com/mark3labs/mcp-go/mcp"
)

func NewChalkApplyTool() mcp.Tool {
	return mcp.NewTool("chalk_apply",
		mcp.WithDescription("Apply Chalk configurations with branch and JSON output. Recommended workflow is to use the chalk_environment tool to set the environment, then use this tool to apply the configuration."),
		mcp.WithString("project_repository",
			mcp.Required(),
			mcp.Description("Path to the root of the Chalk project on disk. Should contain a chalk.yml or chalk.yaml file."),
		),
		mcp.WithString("branch_name",
			mcp.Description("Name of the branch to deploy to. If not provided, defaults to the current git branch."),
		),
	)
}

func ChalkApplyHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectRepository, ok := request.Params.Arguments["project_repository"].(string)
	if !ok {
		return nil, errors.New("project_repository must be a string")
	}

	branchArg := "--branch"
	if branchName, ok := request.Params.Arguments["branch_name"].(string); ok && branchName != "" {
		branchArg += "=" + branchName
	}

	cmd, err := utils.GetChalkCommand(projectRepository, "apply", branchArg)
	if err != nil {
		return nil, errors.Wrap(err, "preparing chalk command")
	}

	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, errors.Wrapf(err, "running chalk command; stdout: %s", out)
	}

	return mcp.NewToolResultText(string(out)), nil
}
