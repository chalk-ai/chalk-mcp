# Chalk MCP Server

An MCP (Model Context Protocol) server that provides tools for interacting with Chalk projects.

## Overview

This MCP server exposes two tools for working with Chalk projects:
- `chalk_features` - Get the list of features from a Chalk project
- `chalk_config` - Get the configuration from a Chalk project

## Requirements

- Go 1.24.0+
- Chalk binary installed on your system

## Installation

```bash
go mod download
go build
```

## Usage

The server runs in stdio mode and expects a Chalk project directory containing either a `chalk.yml` or `chalk.yaml` configuration file.

### Available Tools

#### `chalk_features`
Retrieves the list of features from a Chalk project.

**Parameters:**
- `project_repository` (required): Path to the root of the Chalk project on disk

#### `chalk_config`
Retrieves the configuration from a Chalk project.

**Parameters:**
- `project_repository` (required): Path to the root of the Chalk project on disk

## Testing

Run the test suite:

```bash
go test
```

## License

Licensed under the Apache License, Version 2.0. See [LICENSE](LICENSE) for details.
