use dotenvy::from_filename;
use http_body_util::{BodyExt, Full};
use hyper::{Request, Response};
use hyper::body::{Bytes, Incoming};
use hyper::server::conn::http1;
use hyper_util::rt::TokioIo;
use tokio::net::{TcpListener, TcpStream};
use reqwest::Client;
use std::env;
use std::net::SocketAddr;
use hyper::service::{service_fn};
use std::convert::Infallible;


#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error + Send + Sync>> {
    from_filename("app.env").ok();
    
    let port = env::var("PORT").unwrap_or_else(|_| "3443".to_string()).parse::<u16>().unwrap();
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
    let remote_url = env::var("HOST_SRV_A").expect("HOST_SRV_A must be set");
    let full_url = format!(
        "{}{}",
        remote_url,
        req.uri().path_and_query().map(|x| x.as_str()).unwrap_or("/")
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
}
