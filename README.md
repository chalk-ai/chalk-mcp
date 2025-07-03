# Chalk MCP Server

An MCP (Model Context Protocol) server that provides tools for interacting with Chalk projects.

## Overview

This MCP server exposes tools for working with Chalk projects:
- `chalk_features` - Get the list of features from a Chalk project
- `chalk_config` - Get the configuration from a Chalk project
- `chalk_logs` - Search through Chalk logs using powerful filtering capabilities
- `chalk_environment` - Set or get the Chalk environment for the current project
- `chalk_apply` - Apply Chalk configurations with branch support

## Requirements

- Go 1.24.0+
- [Chalk CLI](https://github.com/chalk-ai/cli) installed on your system

## Installation

Build the MCP server binary:

```bash
go build
```

## Usage

The server runs in stdio mode and expects a Chalk project directory containing either a `chalk.yml` or `chalk.yaml` configuration file.

### Setup

Note the absolute path to your built binary (e.g., `/path/to/chalk-mcp/chalk-mcp`)

### MCP Configuration

#### Claude Code CLI
```bash
claude mcp add-json chalk-mcp '{
  "command": "/path/to/chalk-mcp/chalk-mcp"
}'
```

#### Cursor
For Cursor, open `~/.cursor/mcp.json` and add:
```json
{
   "mcpServers": {
      "chalk-mcp": {
         "command": "/path/to/chalk-mcp/chalk-mcp"
      }
   }
}
```

### Available Tools

#### `chalk_features`
Retrieves the list of features from a Chalk project.

**Parameters:**
- `project_repository` (required): Path to the root of the Chalk project on disk

#### `chalk_config`
Retrieves the configuration from a Chalk project.

**Parameters:**
- `project_repository` (required): Path to the root of the Chalk project on disk

#### `chalk_logs`
Search through Chalk logs using powerful filtering capabilities.

**Parameters:**
- `query` (required): Query to search logs. Examples:
  - `resolver:user_features`
  - `component:engine message:error`
  - `correlation_id:abc-123`
- `project_repository` (required): Path to the root of the Chalk project on disk

#### `chalk_environment`
Set or get the Chalk environment for the current project.

**Parameters:**
- `project_repository` (required): Path to the root of the Chalk project on disk
- `operation` (required): The operation to perform - either `"set"` or `"get"`
- `environment` (optional): If doing a set operation, the environment to set for the current project. Can be an environment ID or name

#### `chalk_apply`
Apply Chalk configurations with branch and JSON output. Recommended workflow is to use the `chalk_environment` tool to set the environment, then use this tool to apply the configuration.

**Parameters:**
- `project_repository` (required): Path to the root of the Chalk project on disk
- `branch_name` (optional): Name of the branch to deploy to. If not provided, defaults to the current git branch

## Testing

Run the test suite:

```bash
go test ./...
```

## License

Licensed under the Apache License, Version 2.0. See [LICENSE](LICENSE) for details.
