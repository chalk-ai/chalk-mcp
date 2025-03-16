use anyhow::Result;
use clap::Parser;
use tokio::signal;
use tracing::{info, error};
use std::net::SocketAddr;
use tonic::transport::Server;

mod proto;
mod service;

use service::DeploymentServiceImpl;
use proto::DeployServiceServer;

#[derive(Parser, Debug)]
#[clap(author, version, about = "Chalk MCP Server")]
struct Cli {
    /// Server host address
    #[clap(long, value_name = "HOST", default_value = "127.0.0.1")]
    host: String,
    
    /// Server port
    #[clap(long, value_name = "PORT", default_value = "50051")]
    port: u16,
}

#[tokio::main]
async fn main() -> Result<()> {
    // Initialize tracing
    tracing_subscriber::fmt::init();
    
    // Parse command line arguments
    let cli = Cli::parse();
    
    // Build the server address
    let addr: SocketAddr = format!("{}:{}", cli.host, cli.port)
        .parse()
        .expect("Failed to parse server address");
    
    info!("Starting Chalk Model Context Protocol server at {}", addr);
    
    // Create the service implementation
    let deployment_service = DeploymentServiceImpl::default();
    
    // Create the server
    let server = Server::builder()
        .add_service(DeployServiceServer::new(deployment_service))
        .serve(addr);
    
    // Handle graceful shutdown with Ctrl+C
    tokio::select! {
        result = server => {
            if let Err(e) = result {
                error!("Server error: {}", e);
            }
        }
        _ = signal::ctrl_c() => {
            info!("Received Ctrl+C, shutting down...");
        }
    }
    
    info!("Server shutdown complete");
    
    Ok(())
}