package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/chalk-ai/chalk-mcp/utils"
	"github.com/cockroachdb/errors"
	"github.com/mark3labs/mcp-go/mcp"
)

func NewChalkQueryTool() mcp.Tool {
	return mcp.NewTool("chalk_query",
		mcp.WithDescription("Query Chalk features using fully qualified feature names (FQNs). Supports input features, output features, and branch specification."),
		mcp.WithString("project_repository",
			mcp.Required(),
			mcp.Description("Path to the root of the Chalk project on disk. Should contain a chalk.yml or chalk.yaml file."),
		),
		mcp.WithObject("input_features",
			mcp.Description("Map of fully qualified feature names (FQNs) to their values for the query. Each key-value pair becomes --in {key}={value}."),
			mcp.AdditionalProperties(map[string]any{"type": "string"}),
		),
		mcp.WithArray("output_features",
			mcp.Description("List of fully qualified feature names (FQNs) to use as outputs for the query."),
			mcp.Items(map[string]any{"type": "string"}),
		),
		mcp.WithString("branch_name",
			mcp.Description("Name of the branch to query against. If not provided, uses the mainline deployment."),
		),
	)
}

func ChalkQueryHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectRepository, ok := request.Params.Arguments["project_repository"].(string)
	if !ok {
		return nil, errors.New("project_repository must be a string")
	}

	var args []string

	if inputFeatures, ok := request.Params.Arguments["input_features"].(map[string]any); ok && len(inputFeatures) > 0 {
		for fqn, rawValue := range inputFeatures {
			args = append(args, "--in", fqn+"="+fmt.Sprintf("%v", rawValue))
		}
	}

	if outputFeatures, ok := request.Params.Arguments["output_features"].([]any); ok && len(outputFeatures) > 0 {
		for _, feature := range outputFeatures {
			if featureStr, ok := feature.(string); ok && featureStr != "" {
				args = append(args, "--out", featureStr)
			}
		}
	}

	if branchName, ok := request.Params.Arguments["branch_name"].(string); ok && branchName != "" {
		args = append(args, "--branch="+branchName)
	}

	finalArgs := append([]string{"query", "--grpc"}, args...)

	cmd, err := utils.GetChalkCommand(projectRepository, finalArgs...)
	if err != nil {
		return nil, errors.Wrap(err, "preparing chalk command")
	}

	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, errors.Wrapf(err, "running chalk command; stdout: %s", out)
	}

	finalArgsStr := strings.Join(finalArgs, " ")

	return mcp.NewToolResultText(fmt.Sprintf("chalk command: %s\nstdout: %s", finalArgsStr, string(out))), nil
}
