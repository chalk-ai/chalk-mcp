package tools

import (
	"context"

	"github.com/chalk-ai/chalk-mcp/utils"
	"github.com/cockroachdb/errors"
	"github.com/mark3labs/mcp-go/mcp"
)

func NewChalkFeaturesTool() mcp.Tool {
	return mcp.NewTool("chalk_features",
		mcp.WithDescription("Get the list of features from a chalk project"),
		mcp.WithString("project_repository",
			mcp.Required(),
			mcp.Description("Path to the root of the Chalk project on disk to fetch features for. Should contain a chalk.yml or chalk.yaml file."),
		),
	)
}

func ChalkFeaturesHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectRepository, ok := request.Params.Arguments["project_repository"].(string)
	if !ok {
		return nil, errors.New("project_repository must be a string")
	}

	cmd, err := utils.GetChalkCommand(projectRepository, "features", "--json")
	if err != nil {
		return nil, errors.Wrap(err, "preparing chalk command")
	}

	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, errors.Wrapf(err, "running chalk command; stdout: %s", string(out))
	}

	return mcp.NewToolResultText(string(out)), nil
}
