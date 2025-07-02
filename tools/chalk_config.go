package tools

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/chalk-ai/chalk-mcp/utils"
	"github.com/mark3labs/mcp-go/mcp"
)

func NewChalkConfigTool() mcp.Tool {
	return mcp.NewTool("chalk_config",
		mcp.WithDescription("Get the chalk config from a chalk project"),
		mcp.WithString("project_repository",
			mcp.Required(),
			mcp.Description("Path to the root of the Chalk project on disk. Should contain a chalk.yml or chalk.yaml file."),
		),
	)
}

func ChalkConfigHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, ok := request.Params.Arguments["project_repository"].(string)
	if !ok {
		return nil, errors.New("project_repository must be a string")
	}

	if name == "" {

	}

	// must be a directory with a chalk.yml or chalk.yaml
	if _, err := os.Stat(name); os.IsNotExist(err) {
		return nil, errors.New("project_repository must exist")
	}
	hasConfig, err := utils.HasChalkConfig(name)
	if err != nil {
		return nil, fmt.Errorf("error checking for chalk config: %w", err)
	}
	if !hasConfig {
		return nil, errors.New("project_repository must contain a chalk.yml or chalk.yaml file")
	}

	// run the 'chalk' program in the directory
	chalkPath, err := utils.FindChalkBinary()
	if err != nil {
		return nil, fmt.Errorf("failed to find chalk binary: %w", err)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home directory: %w", err)
	}

	cmd := exec.Command(chalkPath, "config", "--json")
	cmd.Dir = name
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "XDG_CONFIG_HOME="+homeDir+"/")

	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to run chalk command: %w; stderr: %s", err, out)
	}

	return mcp.NewToolResultText(string(out)), nil
}
