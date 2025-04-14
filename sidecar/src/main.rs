use dotenvy::from_filename;
use http_body_util::{BodyExt, Full};
use hyper::body::{Bytes, Incoming};
use hyper::server::conn::http1;
use hyper::service::service_fn;
use hyper::{Request, Response, Uri};
use hyper_util::rt::TokioIo;
use reqwest::Client;
use std::collections::HashMap;
use std::convert::Infallible;
use std::env;
use std::net::SocketAddr;
use tokio::net::{TcpListener, TcpStream};

use config::Config;

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error + Send + Sync>> {
    dotenvy::from_filename("app.env").ok();

    let port = env::var("PORT")
        .unwrap_or_else(|_| "3443".to_string())
        .parse::<u16>()
        .unwrap();

    let addr = SocketAddr::from(([127, 0, 0, 1], port));

    let listener = TcpListener::bind(addr).await?;

    loop {
        let (stream, _) = listener.accept().await?;

        let io = TokioIo::new(stream);

        tokio::task::spawn(async move {
            if let Err(err) = http1::Builder::new()
                .serve_connection(io, service_fn(proxy_handler))
                .await
            {
                eprintln!("Error serving connection: {:?}", err);
            }
        });
    }
}

pub async fn proxy_handler(req: Request<Incoming>) -> Result<Response<Full<Bytes>>, Infallible> {
    let mut settings = HashMap::<String, String>::new();
    settings.insert(
        "document_service".to_string(),
        "https://localhost:7443".to_string(),
    );
    settings.insert(
        "transaction_service".to_string(),
        "https://localhost:8443".to_string(),
    );

    println!("{}", req.uri());

    let path_and_query = req.uri().to_string();

    let srv = path_and_query.trim_start_matches("/").split("/").next();

    let (srv_path, remote_host) = settings.get_key_value(srv.unwrap()).unwrap();

    let path_and_query = path_and_query.trim_start_matches("/").replace(srv_path, "");

    let remote_url = Uri::builder()
        .authority(remote_host.as_str())
        .path_and_query(path_and_query.as_str())
        .build();

    match remote_url {
        Ok(url) => {
            let full_url = format!(
                "{}{}",
                url,
                req.uri()
                    .path_and_query()
                    .map(|x| x.as_str())
                    .unwrap_or("/")
            );

            let method = req.method().clone();
            let headers = req.headers().clone();

            let collected = req.into_body().collect().await.unwrap();
            let body_bytes = collected.to_bytes();

            let client = Client::builder()
                .danger_accept_invalid_certs(true)
                .build()
                .unwrap();

            let mut request_builder = client.request(method, full_url);
            for (key, value) in headers.iter() {
                if key != "host" {
                    request_builder = request_builder.header(key, value);
                }
            }

            let response = request_builder.body(body_bytes).send().await.unwrap();

            let status = response.status();
            let response_headers = response.headers().clone();
            let body = response.bytes().await.unwrap();

            let mut builder = Response::builder().status(status);
            for (key, value) in response_headers.iter() {
                builder = builder.header(key, value);
            }

            Ok(builder.body(Full::new(Bytes::from(body))).unwrap())
        },
        Err(err) => {
            panic!("{}",err)
        },
    }
}
