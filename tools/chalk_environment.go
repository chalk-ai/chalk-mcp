package tools

import (
	"context"

	"github.com/chalk-ai/chalk-mcp/utils"
	"github.com/cockroachdb/errors"
	"github.com/mark3labs/mcp-go/mcp"
)

type ChalkEnvironmentOperation string

const (
	ChalkEnvironmentOperationSet ChalkEnvironmentOperation = "set"
	ChalkEnvironmentOperationGet ChalkEnvironmentOperation = "get"
)

func NewChalkEnvironmentTool() mcp.Tool {
	return mcp.NewTool("chalk_environment",
		mcp.WithDescription("Set the chalk environment for the current project"),
		mcp.WithString("project_repository",
			mcp.Required(),
			mcp.Description("Path to the root of the Chalk project on disk. Should contain a chalk.yml or chalk.yaml file."),
		),
		mcp.WithString("operation",
			mcp.Required(),
			mcp.Description("The operation to perform on the chalk environment"),
			mcp.Enum(string(ChalkEnvironmentOperationSet), string(ChalkEnvironmentOperationGet)),
		),
		mcp.WithString("environment",
			mcp.Description("If doing a set operation, the environment to set for the current project. Can be an environment ID or name."),
		),
	)
}

func ChalkEnvironmentHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectRepository, ok := request.Params.Arguments["project_repository"].(string)
	if !ok {
		return nil, errors.New("project_repository must be a string")
	}

	operation, ok := request.Params.Arguments["operation"].(string)
	if !ok {
		return nil, errors.New("operation must be a string")
	}

	switch operation {
	case string(ChalkEnvironmentOperationSet):
		environment, ok := request.Params.Arguments["environment"].(string)
		if !ok {
			return nil, errors.New("environment must be a string")
		}

		cmd, err := utils.GetChalkCommand(projectRepository, "environment", environment)
		if err != nil {
			return nil, errors.Wrap(err, "preparing chalk command")
		}

		out, err := cmd.CombinedOutput()
		if err != nil {
			return nil, errors.Wrapf(err, "running chalk command; stdout: %s", out)
		}

		return mcp.NewToolResultText(string(out)), nil
	case string(ChalkEnvironmentOperationGet):
		cmd, err := utils.GetChalkCommand(projectRepository, "environment")
		if err != nil {
			return nil, errors.Wrap(err, "preparing chalk command")
		}

		out, err := cmd.CombinedOutput()
		if err != nil {
			return nil, errors.Wrapf(err, "running chalk command; stdout: %s", out)
		}

		return mcp.NewToolResultText(string(out)), nil
	}

	return nil, errors.Newf("invalid operation: %s", operation)
}
