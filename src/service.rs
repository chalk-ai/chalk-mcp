use async_trait::async_trait;
use tonic::{Request, Response, Status};
use std::time::{SystemTime, UNIX_EPOCH};
use prost_types::Timestamp;

use crate::proto::{
    DeployService, DeploymentStatus, 
    ListDeploymentsRequest, ListDeploymentsResponse,
    GetDeploymentRequest, GetDeploymentResponse,
    GetActiveDeploymentsRequest, GetActiveDeploymentsResponse,
    DeployBranchRequest, DeployBranchResponse,
    CreateBranchFromSourceDeploymentRequest, CreateBranchFromSourceDeploymentResponse,
    SuspendDeploymentRequest, SuspendDeploymentResponse,
    ScaleDeploymentRequest, ScaleDeploymentResponse,
    TagDeploymentRequest, TagDeploymentResponse,
    Deployment,
};

// Error type for our service
#[derive(Debug, thiserror::Error)]
pub enum ServiceError {
    #[error("Deployment not found: {0}")]
    NotFound(String),
    
    #[error("Internal server error: {0}")]
    Internal(String),
    
    #[error("Invalid request: {0}")]
    InvalidRequest(String),
}

impl From<ServiceError> for Status {
    fn from(err: ServiceError) -> Self {
        match err {
            ServiceError::NotFound(msg) => Status::not_found(msg),
            ServiceError::Internal(msg) => Status::internal(msg),
            ServiceError::InvalidRequest(msg) => Status::invalid_argument(msg),
        }
    }
}

// Our deployment service implementation
#[derive(Debug, Default, Clone)]
pub struct DeploymentServiceImpl {
    // In a real implementation, this would probably contain a database client
    // or some other way to store and retrieve deployments
}

#[async_trait]
impl DeployService for DeploymentServiceImpl {
    async fn deploy_branch(
        &self,
        _request: Request<DeployBranchRequest>
    ) -> Result<Response<DeployBranchResponse>, Status> {
        // Not implemented in this example
        Err(Status::unimplemented("Not implemented"))
    }

    async fn create_branch_from_source_deployment(
        &self,
        _request: Request<CreateBranchFromSourceDeploymentRequest>
    ) -> Result<Response<CreateBranchFromSourceDeploymentResponse>, Status> {
        // Not implemented in this example
        Err(Status::unimplemented("Not implemented"))
    }

    async fn list_deployments(
        &self,
        request: Request<ListDeploymentsRequest>,
    ) -> Result<Response<ListDeploymentsResponse>, Status> {
        let req = request.into_inner();
        
        // In a real implementation, we would fetch deployments from a database
        // and handle pagination using the cursor and limit
        
        // For this example, let's return some mock data
        let limit = req.limit.unwrap_or(10) as usize;
        let mock_deployments = generate_mock_deployments(limit);
        
        // If branch_name filter is specified, filter the deployments
        let filtered_deployments = if let Some(branch_name) = req.branch_name {
            mock_deployments
                .into_iter()
                .filter(|d| d.branch == branch_name)
                .collect()
        } else {
            mock_deployments
        };
        
        // Create the response
        let next_cursor = if filtered_deployments.len() == limit {
            // In a real implementation, this would be a token that helps us
            // continue from this point in a subsequent request
            Some("next-page-token".to_string())
        } else {
            None
        };
        
        let response = ListDeploymentsResponse {
            deployments: filtered_deployments,
            cursor: next_cursor,
        };
        
        Ok(Response::new(response))
    }
    
    async fn get_deployment(
        &self,
        request: Request<GetDeploymentRequest>,
    ) -> Result<Response<GetDeploymentResponse>, Status> {
        let req = request.into_inner();
        
        // In a real implementation, we would fetch the deployment from a database
        // For this example, let's generate a mock deployment if the ID looks valid
        if req.deployment_id.is_empty() {
            return Err(Status::invalid_argument("Deployment ID cannot be empty"));
        }
        
        // Generate a mock deployment with the requested ID
        let deployment = generate_mock_deployment(&req.deployment_id);
        
        let response = GetDeploymentResponse {
            deployment: Some(deployment),
        };
        
        Ok(Response::new(response))
    }
    
    async fn get_active_deployments(
        &self,
        _request: Request<GetActiveDeploymentsRequest>,
    ) -> Result<Response<GetActiveDeploymentsResponse>, Status> {
        // In a real implementation, we would fetch active deployments from a database
        // For this example, let's return some mock active deployments
        let active_deployments = generate_mock_active_deployments(5);
        
        let response = GetActiveDeploymentsResponse {
            deployments: active_deployments,
        };
        
        Ok(Response::new(response))
    }
    
    async fn suspend_deployment(
        &self,
        _request: Request<SuspendDeploymentRequest>
    ) -> Result<Response<SuspendDeploymentResponse>, Status> {
        // Not implemented in this example
        Err(Status::unimplemented("Not implemented"))
    }

    async fn scale_deployment(
        &self,
        _request: Request<ScaleDeploymentRequest>
    ) -> Result<Response<ScaleDeploymentResponse>, Status> {
        // Not implemented in this example
        Err(Status::unimplemented("Not implemented"))
    }

    async fn tag_deployment(
        &self,
        _request: Request<TagDeploymentRequest>
    ) -> Result<Response<TagDeploymentResponse>, Status> {
        // Not implemented in this example
        Err(Status::unimplemented("Not implemented"))
    }
}

// Helper function to generate mock deployments (for demo purposes)
fn generate_mock_deployments(count: usize) -> Vec<Deployment> {
    (0..count)
        .map(|i| {
            let id = format!("deployment-{}", i);
            generate_mock_deployment(&id)
        })
        .collect()
}

// Helper function to generate a single mock deployment (for demo purposes)
fn generate_mock_deployment(id: &str) -> Deployment {
    // Get current time as a timestamp
    let now = SystemTime::now()
        .duration_since(UNIX_EPOCH)
        .unwrap_or_default();
    
    let timestamp = Timestamp {
        seconds: now.as_secs() as i64,
        nanos: now.subsec_nanos() as i32,
    };
    
    // Generate a mock deployment
    Deployment {
        id: id.to_string(),
        environment_id: "production".to_string(),
        status: DeploymentStatus::Success as i32,
        deployment_tags: vec!["latest".to_string(), "stable".to_string()],
        cloud_build_id: format!("build-{}", id),
        triggered_by: "user@example.com".to_string(),
        requirements_filepath: Some("requirements.txt".to_string()),
        dockerfile_filepath: Some("Dockerfile".to_string()),
        runtime: Some("python3.8".to_string()),
        chalkpy_version: "0.5.0".to_string(),
        raw_dependency_hash: "abc123".to_string(),
        final_dependency_hash: Some("def456".to_string()),
        is_preview_deployment: Some(false),
        created_at: Some(timestamp.clone()),
        updated_at: Some(timestamp.clone()),
        git_commit: "a1b2c3d4e5f6".to_string(),
        git_pr: "42".to_string(),
        git_branch: "main".to_string(),
        git_author_email: "author@example.com".to_string(),
        branch: "main".to_string(),
        project_settings: "{}".to_string(),
        requirements_files: Some("[]".to_string()),
        git_tag: "v1.0.0".to_string(),
        base_image_sha: "sha256:12345".to_string(),
        status_changed_at: Some(timestamp),
        pinned_platform_version: None,
        preview_deployment_tag: None,
    }
}

// Helper function to generate mock active deployments (for demo purposes)
fn generate_mock_active_deployments(count: usize) -> Vec<Deployment> {
    (0..count)
        .map(|i| {
            let id = format!("active-deployment-{}", i);
            let mut deployment = generate_mock_deployment(&id);
            deployment.status = DeploymentStatus::Success as i32;
            deployment
        })
        .collect()
}