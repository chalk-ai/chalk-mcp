package tools

import (
	"context"
	"errors"
	"fmt"

	"github.com/chalk-ai/chalk-mcp/utils"
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
		return nil, fmt.Errorf("failed to prepare chalk command: %w", err)
	}

	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to run chalk command: %w; stdout: %s", err, out)
	}

	return mcp.NewToolResultText(string(out)), nil
}
