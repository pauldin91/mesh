use hyper::body::Incoming;
use hyper::server::conn::http1;
use hyper::service::service_fn;
use hyper::Request;
use hyper_util::rt::TokioIo;
use std::env;
use std::net::SocketAddr;
use std::sync::Arc;
use tokio::net::TcpListener;
use proxy::Handler;

mod proxy;
mod utils;

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error + Send + Sync>> {
    dotenvy::from_filename("app.env").ok();

    let port = env::var("PORT")
        .unwrap_or_else(|_| "3443".to_string())
        .parse::<u16>()
        .unwrap();

    let addr = SocketAddr::from(([127, 0, 0, 1], port));

    let listener = TcpListener::bind(addr).await?;

    let handler = Arc::new(Handler::new("config.yml"));

    loop {
        let (stream, _) = listener.accept().await?;

        let io = TokioIo::new(stream);

        let handler_clone = Arc::clone(&handler); 

        tokio::task::spawn(async move {
            let service = service_fn(move |req: Request<Incoming>| {
                let handler = Arc::clone(&handler_clone);
                async move { handler.proxy_handler(req).await }
            });
        
            if let Err(err) = http1::Builder::new()
                .serve_connection(io, service)
                .await
            {
                eprintln!("Error serving connection: {:?}", err);
            }
        });
    }
}
