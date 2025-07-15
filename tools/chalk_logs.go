package tools

import (
	"context"
	"errors"
	"fmt"

	"github.com/chalk-ai/chalk-mcp/utils"
	"github.com/mark3labs/mcp-go/mcp"
)

func NewChalkLogsTool() mcp.Tool {
	return mcp.NewTool("chalk_logs",
		mcp.WithDescription("Search through Chalk logs using powerful filtering capabilities"),
		mcp.WithString("query",
			mcp.Required(),
			mcp.Description("Query to search logs. Examples: 'resolver:user_features', 'component:engine message:error', 'correlation_id:abc-123'"),
		),
		mcp.WithString("project_repository",
			mcp.Required(),
			mcp.Description("Path to the root of the Chalk project on disk. Should contain a chalk.yml or chalk.yaml file."),
		),
	)
}

func ChalkLogsHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectRepository, ok := request.Params.Arguments["project_repository"].(string)
	if !ok {
		return nil, errors.New("project_repository must be a string")
	}

	query, ok := request.Params.Arguments["query"].(string)
	if !ok {
		return nil, errors.New("query must be a string")
	}

	cmd, err := utils.GetChalkCommand(ctx, projectRepository, "logs", "--query", query)
	if err != nil {
		return nil, fmt.Errorf("preparing chalk command: %w", err)
	}

	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("running chalk logs command: %w; output: %s", err, out)
	}

	return mcp.NewToolResultText(string(out)), nil
}
