pub mod service_config; // make the submodule available

pub use service_config::*; // re-export contents so `utils::foo` works
