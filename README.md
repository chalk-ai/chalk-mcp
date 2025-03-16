# Chalk Model Context Protocol Server

A Rust implementation of a model context protocol server for the Chalk API. This server provides gRPC services to interact with Chalk deployments.

## Features

- **ListDeployments**: Query and filter deployments with pagination support
- **GetDeployment**: Retrieve detailed information about a specific deployment
- **GetActiveDeployments**: Fetch currently active deployments

## Requirements

- Rust 1.56 or higher
- Access to Chalk API protobuf definitions

## Building

```bash
cargo build
```

## Running

```bash
# Run with default settings (listens on 127.0.0.1:50051)
cargo run

# Run with custom host and port
cargo run -- --host 0.0.0.0 --port 8080

# Run with a configuration file
cargo run -- --config config.toml
```

## Configuration

The server can be configured using a TOML configuration file:

```toml
[server]
host = "0.0.0.0"
port = 8080
```

Environment variables can also be used to override settings:

```bash
CHALK_MCP_SERVER_HOST=0.0.0.0 CHALK_MCP_SERVER_PORT=8080 cargo run
```

## Architecture

- **proto**: Generated Protobuf code for Chalk API
- **services**: Implementation of gRPC services
- **server**: gRPC server setup and configuration
- **config**: Application configuration management

## Implementation

The server is implemented using the following Rust crates:
- `tonic`: gRPC framework
- `tokio`: Async runtime
- `clap`: Command-line argument parsing
- `config`: Configuration management
- `tracing`: Logging and instrumentation