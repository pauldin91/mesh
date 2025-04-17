use config::{Config, ConfigError};
use std::collections::HashMap;
use serde::Deserialize;


#[derive(Debug, Deserialize, Clone)]
pub struct ServiceConfig {
    pub services: HashMap<String, String>,
}


impl ServiceConfig {
    
    pub fn load_from_file(path: &str) -> Self {
        Config::builder()
            .add_source(config::File::with_name(path))
            .build()
            .unwrap()
            .try_deserialize()
            .unwrap()
    }

    pub fn get_settings(&self) -> HashMap<String,String>{
        self.services.clone()
    }
}
