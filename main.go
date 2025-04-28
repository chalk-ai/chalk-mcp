package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"os"
	"os/exec"
)

func main() {
	// Create MCP server
	s := server.NewMCPServer(
		"Demo 🚀",
		"1.0.0",
	)

	// Add tool
	tool := mcp.NewTool("chalk_features",
		mcp.WithDescription("Get the list of features from a chalk project"),
		mcp.WithString("project_repository",
			mcp.Required(),
			mcp.Description("Path to the root of the Chalk project on disk to fetch features for. Should contain a chalk.yml file."),
		),
	)

	tool2 := mcp.NewTool("chalk_config",
		mcp.WithDescription("Get the chalk config from a chalk project"),
		mcp.WithString("project_repository",
			mcp.Required(),
			mcp.Description("Path to the root of the Chalk project on disk. Should contain a chalk.yml file."),
		),
	)

	// Add tool handler
	s.AddTool(tool, helloHandler)
	s.AddTool(tool2, configHandler)

	// Start the stdio server
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}

func helloHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, ok := request.Params.Arguments["project_repository"].(string)
	if !ok {
		return nil, errors.New("project_repository must be a string")
	}

	if name == "" {

	}

	// must be a directory
	if _, err := os.Stat(name); os.IsNotExist(err) {
		return nil, errors.New("project_repository must exist")
	}

	// must be a directory with a chalk.yml
	if _, err := os.Stat(name + "/chalk.yml"); os.IsNotExist(err) {
		return nil, errors.New("project_repository must contain a chalk.yml file")
	}

	// run the 'chalk' program in the directory

	cmd := exec.Command("/Users/andrew/.chalk/bin/chalk-latest", "features", "--json")
	cmd.Dir = name
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "XDG_CONFIG_HOME=/Users/andrew/")
	out, err := cmd.CombinedOutput()

	// return the json output

	if err != nil {

		return nil, fmt.Errorf("failed to run chalk command: %w; stderr: %s", err, out)
	}

	// parse the json output

	return mcp.NewToolResultText(string(out)), nil
}

func configHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, ok := request.Params.Arguments["project_repository"].(string)
	if !ok {
		return nil, errors.New("project_repository must be a string")
	}

	if name == "" {

	}

	// must be a directory
	if _, err := os.Stat(name); os.IsNotExist(err) {
		return nil, errors.New("project_repository must exist")
	}

	// must be a directory with a chalk.yml
	if _, err := os.Stat(name + "/chalk.yml"); os.IsNotExist(err) {
		return nil, errors.New("project_repository must contain a chalk.yml file")
	}

	// run the 'chalk' program in the directory

	cmd := exec.Command("/Users/andrew/.chalk/bin/chalk-latest", "config", "--json")
	cmd.Dir = name
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "XDG_CONFIG_HOME=/Users/andrew/")

	out, err := cmd.CombinedOutput()

	// return the json output

	if err != nil {

		return nil, fmt.Errorf("failed to run chalk command: %w; stderr: %s", err, out)
	}

	// parse the json output

	return mcp.NewToolResultText(string(out)), nil
}
