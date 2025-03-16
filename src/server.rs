use std::net::SocketAddr;
use tonic::transport::Server;
use anyhow::Result;
use tracing::info;

use crate::proto::chalk::server::v1::deploy_service_server::DeployServiceServer;
use crate::services::DeploymentServiceImpl;

/// Server configuration parameters
#[derive(Debug, Clone, serde::Serialize, serde::Deserialize)]
pub struct ServerConfig {
    pub host: String,
    pub port: u16,
}

impl Default for ServerConfig {
    fn default() -> Self {
        Self {
            host: "127.0.0.1".to_string(),
            port: 50051,
        }
    }
}

/// GrpcServer manages the gRPC server instance
pub struct GrpcServer {
    config: ServerConfig,
}

impl GrpcServer {
    /// Create a new GrpcServer with the given config
    pub fn new(config: ServerConfig) -> Self {
        Self { config }
    }
    
    /// Start the gRPC server
    pub async fn start(&self) -> Result<()> {
        // Create service implementations
        let deployment_service = DeploymentServiceImpl::default();
        
        // Build the server address
        let addr: SocketAddr = format!("{}:{}", self.config.host, self.config.port)
            .parse()
            .expect("Failed to parse server address");
        
        // Start the server
        info!("Starting gRPC server at {}", addr);
        
        Server::builder()
            .add_service(DeployServiceServer::new(deployment_service))
            // Add more services here as needed
            .serve(addr)
            .await?;
        
        Ok(())
    }
}