package main

import (
	"fmt"

	"github.com/chalk-ai/chalk-mcp/tools"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	s := server.NewMCPServer(
		"Chalk MCP 🚀",
		"1.0.0",
	)

	chalkTools := []tools.Tool{
		tools.NewChalkConfigTool(nil),
		tools.NewChalkFeaturesTool(nil),
		tools.NewChalkEnvironmentTool(nil),
		tools.NewChalkApplyTool(nil),
		tools.NewChalkQueryTool(nil),
		tools.NewChalkLogsTool(nil),
		tools.NewChalkLintTool(nil),
	}

	for _, tool := range chalkTools {
		s.AddTool(tools.GenerateMetadata(tool), tools.CreateHandler(tool))
	}

	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
