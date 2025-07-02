package main

import (
	"fmt"

	"github.com/chalk-ai/chalk-mcp/tools"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	s := server.NewMCPServer(
		"Demo 🚀",
		"1.0.0",
	)

	s.AddTool(tools.NewChalkConfigTool(), tools.ChalkConfigHandler)
	s.AddTool(tools.NewChalkFeaturesTool(), tools.ChalkFeaturesHandler)

	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
