use crate::utils::{self, ServiceConfig};
use http_body_util::{BodyExt, Full};
use hyper::body::{Bytes, Incoming};
use hyper::{Request, Response, StatusCode};
use reqwest::Client;
use std::convert::Infallible;
use std::sync::Arc; 

#[derive(Clone)]
pub struct Handler{
    config : Arc<ServiceConfig>,
    client: Client
}


impl Handler {

    pub fn new(cfg_path: &str)-> Self{
        let cfg = Arc::new(ServiceConfig::load_from_file(cfg_path));
        
        let client = Client::builder()
            .danger_accept_invalid_certs(true)
            .build()
            .unwrap();

        Self{config:cfg,client}
    }

    pub async fn proxy_handler(&self,req: Request<Incoming>) -> Result<Response<Full<Bytes>>, Infallible> {
        let method = req.method().clone();
        let headers = req.headers().clone();
        let uri_path = req.uri().path();

        let service_key = uri_path
            .trim_start_matches('/')
            .split('/')
            .next()
            .unwrap_or("");

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

        let path_to_forward = uri_path
            .strip_prefix(&format!("/{}", service_key))
            .unwrap_or("");
        let forward_path = if path_to_forward.is_empty() {
            "/"
        } else {
            path_to_forward
        };

        let full_url = format!("{}/{}", base_url, forward_path.trim_start_matches('/'));

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


        let mut request_builder = self.client.request(method, full_url);
        for (key, value) in headers.iter() {
            if key != "host" {
                request_builder = request_builder.header(key, value);
            }
        }

        let response = match request_builder.body(body_bytes).send().await {
            Ok(resp) => resp,
            Err(e) => {
                return Ok(Response::builder()
                    .status(StatusCode::BAD_GATEWAY)
                    .body(Full::new(Bytes::from(format!(
                        "Upstream request failed: {}",
                        e
                    ))))
                    .unwrap());
            }
        };

        let status = response.status();
        let resp_headers = response.headers().clone();
        let resp_body = match response.bytes().await {
            Ok(body) => body,
            Err(_) => Bytes::new(),
        };

        let mut builder = Response::builder().status(status);
        for (key, value) in resp_headers.iter() {
            builder = builder.header(key, value);
        }

        Ok(builder.body(Full::new(resp_body)).unwrap())
    }
}
