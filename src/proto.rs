// This file contains manually defined Protobuf types for the Chalk API
use tonic::{Request, Response, Status};
use tonic::transport::Server;
use tonic::server::NamedService;
use prost_types::Timestamp;

/// Deployment status enum
#[derive(Clone, Copy, Debug, PartialEq, Eq, Hash, PartialOrd, Ord, Default)]
#[repr(i32)]
pub enum DeploymentStatus {
    #[default]
    Unspecified = 0,
    Unknown = 1,
    Pending = 2,
    Queued = 3,
    Working = 4,
    Success = 5,
    Failure = 6,
    InternalError = 7,
    Timeout = 8,
    Cancelled = 9,
    Expired = 10,
    BootErrors = 11,
    AwaitingSource = 12,
    Deploying = 13,
}

impl From<i32> for DeploymentStatus {
    fn from(value: i32) -> Self {
        match value {
            0 => DeploymentStatus::Unspecified,
            1 => DeploymentStatus::Unknown,
            2 => DeploymentStatus::Pending,
            3 => DeploymentStatus::Queued,
            4 => DeploymentStatus::Working,
            5 => DeploymentStatus::Success,
            6 => DeploymentStatus::Failure,
            7 => DeploymentStatus::InternalError,
            8 => DeploymentStatus::Timeout,
            9 => DeploymentStatus::Cancelled,
            10 => DeploymentStatus::Expired,
            11 => DeploymentStatus::BootErrors,
            12 => DeploymentStatus::AwaitingSource,
            13 => DeploymentStatus::Deploying,
            _ => DeploymentStatus::Unspecified,
        }
    }
}

/// Instance sizing for deployments
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct InstanceSizing {
    #[prost(uint32, optional, tag="1")]
    pub min_instances: Option<u32>,
    #[prost(uint32, optional, tag="2")]
    pub max_instances: Option<u32>,
}

/// Deployment model
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct Deployment {
    #[prost(string, tag="1")]
    pub id: String,
    #[prost(string, tag="2")]
    pub environment_id: String,
    #[prost(enumeration="DeploymentStatus", tag="3")]
    pub status: i32,
    #[prost(string, repeated, tag="4")]
    pub deployment_tags: ::prost::alloc::vec::Vec<String>,
    #[prost(string, tag="5")]
    pub cloud_build_id: String,
    #[prost(string, tag="6")]
    pub triggered_by: String,
    #[prost(string, optional, tag="7")]
    pub requirements_filepath: Option<String>,
    #[prost(string, optional, tag="8")]
    pub dockerfile_filepath: Option<String>,
    #[prost(string, optional, tag="9")]
    pub runtime: Option<String>,
    #[prost(string, tag="10")]
    pub chalkpy_version: String,
    #[prost(string, tag="11")]
    pub raw_dependency_hash: String,
    #[prost(string, optional, tag="12")]
    pub final_dependency_hash: Option<String>,
    #[prost(bool, optional, tag="13")]
    pub is_preview_deployment: Option<bool>,
    #[prost(message, optional, tag="14")]
    pub created_at: Option<Timestamp>,
    #[prost(message, optional, tag="15")]
    pub updated_at: Option<Timestamp>,
    #[prost(string, tag="16")]
    pub git_commit: String,
    #[prost(string, tag="17")]
    pub git_pr: String,
    #[prost(string, tag="18")]
    pub git_branch: String,
    #[prost(string, tag="19")]
    pub git_author_email: String,
    #[prost(string, tag="20")]
    pub branch: String,
    #[prost(string, tag="21")]
    pub project_settings: String,
    #[prost(string, optional, tag="22")]
    pub requirements_files: Option<String>,
    #[prost(string, tag="23")]
    pub git_tag: String,
    #[prost(string, tag="24")]
    pub base_image_sha: String,
    #[prost(message, optional, tag="25")]
    pub status_changed_at: Option<Timestamp>,
    #[prost(string, optional, tag="26")]
    pub pinned_platform_version: Option<String>,
    #[prost(string, optional, tag="27")]
    pub preview_deployment_tag: Option<String>,
}

/// Request and response types for ListDeployments
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct ListDeploymentsRequest {
    #[prost(string, optional, tag="1")]
    pub cursor: Option<String>,
    #[prost(int32, optional, tag="2")]
    pub limit: Option<i32>,
    #[prost(bool, optional, tag="3")]
    pub include_branch: Option<bool>,
    #[prost(string, optional, tag="4")]
    pub branch_name: Option<String>,
}

#[derive(Clone, PartialEq, ::prost::Message)]
pub struct ListDeploymentsResponse {
    #[prost(message, repeated, tag="1")]
    pub deployments: ::prost::alloc::vec::Vec<Deployment>,
    #[prost(string, optional, tag="2")]
    pub cursor: Option<String>,
}

/// Request and response types for GetDeployment
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GetDeploymentRequest {
    #[prost(string, tag="1")]
    pub deployment_id: String,
}

#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GetDeploymentResponse {
    #[prost(message, optional, tag="1")]
    pub deployment: Option<Deployment>,
}

/// Request and response types for GetActiveDeployments
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GetActiveDeploymentsRequest {
    // Empty request
}

#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GetActiveDeploymentsResponse {
    #[prost(message, repeated, tag="1")]
    pub deployments: ::prost::alloc::vec::Vec<Deployment>,
}

/// Request and response types for other operations (stubs for completeness)
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct DeployBranchRequest {
    #[prost(string, tag="1")]
    pub branch_name: String,
    #[prost(bool, tag="2")]
    pub reset_branch: bool,
    #[prost(bytes, tag="3")]
    pub archive: ::prost::alloc::vec::Vec<u8>,
    #[prost(bool, tag="4")]
    pub is_hot_deploy: bool,
}

#[derive(Clone, PartialEq, ::prost::Message)]
pub struct DeployBranchResponse {
    #[prost(string, tag="1")]
    pub deployment_id: String,
}

#[derive(Clone, PartialEq, ::prost::Message)]
pub struct CreateBranchFromSourceDeploymentRequest {
    #[prost(string, tag="1")]
    pub branch_name: String,
    #[prost(string, tag="2")]
    pub source_branch_name: String,
}

#[derive(Clone, PartialEq, ::prost::Message)]
pub struct CreateBranchFromSourceDeploymentResponse {
    #[prost(string, tag="1")]
    pub deployment_id: String,
    #[prost(bool, tag="4")]
    pub branch_already_exists: bool,
}

#[derive(Clone, PartialEq, ::prost::Message)]
pub struct SuspendDeploymentRequest {
    #[prost(string, tag="1")]
    pub deployment_id: String,
}

#[derive(Clone, PartialEq, ::prost::Message)]
pub struct SuspendDeploymentResponse {
    #[prost(message, optional, tag="1")]
    pub deployment: Option<Deployment>,
}

#[derive(Clone, PartialEq, ::prost::Message)]
pub struct ScaleDeploymentRequest {
    #[prost(string, tag="1")]
    pub deployment_id: String,
    #[prost(message, optional, tag="2")]
    pub sizing: Option<InstanceSizing>,
}

#[derive(Clone, PartialEq, ::prost::Message)]
pub struct ScaleDeploymentResponse {
    #[prost(message, optional, tag="1")]
    pub deployment: Option<Deployment>,
}

#[derive(Clone, PartialEq, ::prost::Message)]
pub struct TagDeploymentRequest {
    #[prost(string, tag="1")]
    pub deployment_id: String,
    #[prost(string, tag="2")]
    pub tag: String,
}

#[derive(Clone, PartialEq, ::prost::Message)]
pub struct TagDeploymentResponse {
    #[prost(message, optional, tag="1")]
    pub deployment: Option<Deployment>,
    #[prost(string, optional, tag="2")]
    pub untagged_deployment_id: Option<String>,
}

/// Service trait for the DeployService
#[tonic::async_trait]
pub trait DeployService: Send + Sync + 'static {
    async fn deploy_branch(
        &self,
        request: Request<DeployBranchRequest>,
    ) -> Result<Response<DeployBranchResponse>, Status>;

    async fn create_branch_from_source_deployment(
        &self,
        request: Request<CreateBranchFromSourceDeploymentRequest>,
    ) -> Result<Response<CreateBranchFromSourceDeploymentResponse>, Status>;

    async fn get_deployment(
        &self,
        request: Request<GetDeploymentRequest>,
    ) -> Result<Response<GetDeploymentResponse>, Status>;

    async fn list_deployments(
        &self,
        request: Request<ListDeploymentsRequest>,
    ) -> Result<Response<ListDeploymentsResponse>, Status>;

    async fn get_active_deployments(
        &self,
        request: Request<GetActiveDeploymentsRequest>,
    ) -> Result<Response<GetActiveDeploymentsResponse>, Status>;

    async fn suspend_deployment(
        &self,
        request: Request<SuspendDeploymentRequest>,
    ) -> Result<Response<SuspendDeploymentResponse>, Status>;

    async fn scale_deployment(
        &self,
        request: Request<ScaleDeploymentRequest>,
    ) -> Result<Response<ScaleDeploymentResponse>, Status>;

    async fn tag_deployment(
        &self,
        request: Request<TagDeploymentRequest>,
    ) -> Result<Response<TagDeploymentResponse>, Status>;
}

// Simple server implementation
#[derive(Debug, Clone)]
pub struct DeployServiceServer<T> {
    inner: T,
}

impl<T> DeployServiceServer<T>
where
    T: DeployService,
{
    pub fn new(inner: T) -> Self {
        Self { inner }
    }
}

impl<T: DeployService + Clone + 'static> NamedService for DeployServiceServer<T> {
    const NAME: &'static str = "chalk.server.v1.DeployService";
}

impl<T: DeployService + Clone + 'static> tower::Service<http::Request<hyper::Body>> for DeployServiceServer<T> {
    type Response = http::Response<tonic::body::BoxBody>;
    type Error = std::convert::Infallible;
    type Future = futures_util::future::BoxFuture<'static, Result<Self::Response, Self::Error>>;

    fn poll_ready(&mut self, _cx: &mut std::task::Context<'_>) -> std::task::Poll<Result<(), Self::Error>> {
        std::task::Poll::Ready(Ok(()))
    }

    fn call(&mut self, req: http::Request<hyper::Body>) -> Self::Future {
        // Just return unimplemented for simplicity in this example
        Box::pin(async move {
            let svc = tonic::Status::unimplemented("Not implemented");
            Ok(http::Response::builder()
                .status(200)
                .header("grpc-status", "12")
                .header("content-type", "application/grpc")
                .body(tonic::body::empty_body())
                .unwrap())
        })
    }
}