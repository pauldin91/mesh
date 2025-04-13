use dotenvy::dotenv;
use std::env;

#[tokio::main]
async fn main() {
    dotenv().ok();
    let address = env::var("PORT").unwrap_or_else(|_|"4443".to_string());

    let server = Server::bind(&address).serve();
}
