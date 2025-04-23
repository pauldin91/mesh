use crate::utils::{self, ServiceConfig};
use http_body_util::{BodyExt, Full};
use hyper::{
    body::{Bytes, Incoming},
    client::conn,
    Request, Response, StatusCode, Uri,
};
use hyper_util::rt::TokioIo;
use std::{convert::Infallible, sync::Arc};
use tokio::net::TcpStream;

#[derive(Clone)]
pub struct Handler {
    config: Arc<ServiceConfig>,
}

impl Handler {
    pub fn new(cfg_path: &str) -> Self {
        let cfg = Arc::new(ServiceConfig::load_from_file(cfg_path));
        Self { config: cfg }
    }

    pub async fn proxy_handler(
        &self,
        req: Request<Incoming>,
    ) -> Result<Response<Full<Bytes>>, Infallible> {
        let method = req.method().clone();
        let headers = req.headers().clone();
        let uri_path = req.uri().path();

        // Check for Swagger paths, make sure it's properly forwarded
        let is_swagger_request = uri_path.starts_with("/swagger");

        // For swagger paths, strip the leading "/swagger" and forward
        let service_key = if is_swagger_request {
            "swagger"
        } else {
            uri_path
                .trim_start_matches('/')
                .split('/')
                .next()
                .unwrap_or("")
        };

        let base_url = match self.config.services.get(service_key) {
            Some(url) => url.trim_end_matches('/'),
            None => {
                return Ok(Response::builder()
                    .status(StatusCode::NOT_FOUND)
                    .body(Full::new(Bytes::from(format!(
                        "Unknown service: {}",
                        service_key
                    ))))
                    .unwrap());
            }
        };

        // For Swagger, remove "/swagger" prefix to forward requests correctly
        let path_to_forward = if is_swagger_request {
            uri_path.strip_prefix("/swagger").unwrap_or("")
        } else {
            &uri_path.replacen(&format!("/{}", service_key), "", 1)
        };

        let forward_path = if path_to_forward.is_empty() {
            "/"
        } else {
            path_to_forward
        };

        let full_url = format!("{}/{}", base_url, forward_path.trim_start_matches('/'));
        let uri: Uri = match full_url.parse() {
            Ok(uri) => uri,
            Err(_) => {
                return Ok(Response::builder()
                    .status(StatusCode::BAD_REQUEST)
                    .body(Full::new(Bytes::from("Invalid URI")))
                    .unwrap());
            }
        };

        // Collect the request body (if any)
        let collected = match req.into_body().collect().await {
            Ok(body) => body,
            Err(_) => {
                return Ok(Response::builder()
                    .status(StatusCode::BAD_REQUEST)
                    .body(Full::new(Bytes::from("Failed to read request body")))
                    .unwrap());
            }
        };

        let body_bytes = collected.to_bytes();

        // Get authority to connect to the backend service
        let authority = match uri.authority() {
            Some(a) => a.clone(),
            None => {
                return Ok(Response::builder()
                    .status(StatusCode::BAD_REQUEST)
                    .body(Full::new(Bytes::from("Missing authority in URI")))
                    .unwrap());
            }
        };

        // Create the TCP connection to backend service
        let stream = match TcpStream::connect(authority.to_string()).await {
            Ok(s) => s,
            Err(_) => {
                return Ok(Response::builder()
                    .status(StatusCode::BAD_GATEWAY)
                    .body(Full::new(Bytes::from("Failed to connect to backend")))
                    .unwrap());
            }
        };

        let io = TokioIo::new(stream);

        // Handshake to establish connection
        let (mut sender, conn) = match conn::http1::handshake(io).await {
            Ok(parts) => parts,
            Err(_) => {
                return Ok(Response::builder()
                    .status(StatusCode::INTERNAL_SERVER_ERROR)
                    .body(Full::new(Bytes::from("Handshake failed")))
                    .unwrap());
            }
        };

        tokio::spawn(async move {
            if let Err(e) = conn.await {
                eprintln!("Connection error: {:?}", e);
            }
        });

        // Build the proxied request
        let mut req_builder = Request::builder().method(method).uri(uri.clone());

        for (key, value) in headers {
            req_builder = req_builder.header(key.unwrap(), value);
        }

        let proxied_request = req_builder
            .body(Full::new(body_bytes))
            .expect("Failed to build proxied request");

        // Send the proxied request
        let response = match sender.send_request(proxied_request).await {
            Ok(res) => res,
            Err(_) => {
                return Ok(Response::builder()
                    .status(StatusCode::BAD_GATEWAY)
                    .body(Full::new(Bytes::from("Request failed")))
                    .unwrap());
            }
        };

        // Read the response from the backend
        let status = response.status();
        let resp_headers = response.headers().clone();
        let collected = match response.into_body().collect().await {
            Ok(body) => body,
            Err(_) => {
                return Ok(Response::builder()
                    .status(StatusCode::BAD_GATEWAY)
                    .body(Full::new(Bytes::from("Failed to read response body")))
                    .unwrap());
            }
        };

        let body_bytes = collected.to_bytes();

        // Build and return the final response
        let mut builder = Response::builder().status(status);
        for (key, value) in resp_headers.iter() {
            builder = builder.header(key, value);
        }

        Ok(builder.body(Full::new(body_bytes)).unwrap())
    }
}
