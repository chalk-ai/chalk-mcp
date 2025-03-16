use anyhow::Result;
use serde::{Deserialize, Serialize};
use config::{Config, Environment, File};
use std::path::Path;

use crate::server::ServerConfig;

/// Application configuration
#[derive(Debug, Deserialize, Serialize, Clone)]
pub struct AppConfig {
    pub server: ServerConfig,
    // Add more configuration sections as needed
}

impl Default for AppConfig {
    fn default() -> Self {
        Self {
            server: ServerConfig::default(),
        }
    }
}

impl AppConfig {
    /// Load configuration from files and environment variables
    pub fn load<P: AsRef<Path>>(config_path: Option<P>) -> Result<Self> {
        let mut builder = Config::builder();
        
        // Start with default config
        builder = builder.add_source(config::Config::try_from(&AppConfig::default())?);
        
        // Load from config files if provided
        if let Some(path) = config_path {
            if path.as_ref().exists() {
                builder = builder.add_source(File::from(path.as_ref()));
            }
        }
        
        // Override with environment variables prefixed with "CHALK_MCP_"
        builder = builder.add_source(
            Environment::with_prefix("CHALK_MCP").separator("_")
        );
        
        let config = builder.build()?;
        let app_config = config.try_deserialize::<AppConfig>()?;
        
        Ok(app_config)
    }
}