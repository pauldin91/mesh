use dotenvy::dotenv;
use http_body_util::{BodyExt, Full};
use hyper::body::{Body,Bytes, Incoming};
use hyper::server::conn::http1;
use hyper::{Request, Response, Uri};
use hyper_util::rt::TokioIo;
use tokio::net::{TcpListener, TcpStream};
use std::env;
use std::fmt::format;
use hyper::client::conn::http1::Builder;
use std::net::SocketAddr;
use hyper::service::{ service_fn,};
use std::convert::Infallible;



#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error + Send + Sync>>{
    dotenv().ok();
    let port = env::var("PORT").unwrap_or_else(|_|"4443".to_string()).parse::<u16>().unwrap();
    let addr = SocketAddr::from(([127, 0, 0, 1],port ));
    

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

async fn proxy_handler(req: Request<Incoming>) -> Result<Response<Full<Bytes>>, Infallible> {
    let remote_url = env::var("HOST_SRV_A").unwrap(); 
    let base_uri: Uri = remote_url.parse().expect("HOST_SRV_A must be a valid URI");

    let path_and_query = req.uri().path_and_query().map(|pq| pq.as_str()).unwrap_or("/");
    let full_uri = format!("{}{}", base_uri, path_and_query);

    let authority = base_uri.authority().unwrap().to_string(); 
    let stream = TcpStream::connect(authority).await.unwrap();
    let io = TokioIo::new(stream);

    let (mut client, conn) = Builder::new()
        .preserve_header_case(true)
        .title_case_headers(true)
        .handshake(io)
        .await
        .unwrap();

    tokio::spawn(async move {
        if let Err(err) = conn.await {
            eprintln!("Connection error: {:?}", err);
        }
    });

    let (parts, body) = req.into_parts();
    let proxied_request = Request::from_parts(parts, body);
    let proxied_request = proxied_request.map(|b| b.boxed());

    let response = client.send_request(proxied_request).await.unwrap();

    let (parts, body) = response.into_parts();
    let body_bytes = body.collect().await.unwrap().to_bytes();
    let full_response = Response::from_parts(parts, Full::new(body_bytes));

    Ok(full_response)
}


