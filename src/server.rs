use tokio::net::{TcpListener, TcpStream};

#[derive(Debug)]
pub struct Server {
    port: i32,
    detach: bool,
    debug: bool,
}

#[derive(Debug)]
struct Handler {
    stream: TcpStream,
}

impl Handler {
    pub fn new(stream: TcpStream) -> Self {
        Handler { stream }
    }
}

impl Server {
    pub fn new(port: i32, detach: bool, debug: bool) -> Self {
        Server {
            port,
            detach,
            debug,
        }
    }

    pub async fn start(&self) {
        println!("{:#?}", self);
        let addr = format!("0.0.0.0:{}", self.port);
        let listener = TcpListener::bind(addr).await.unwrap();

        loop {
            let (stream, _) = listener.accept().await.unwrap();
            tokio::spawn(async move { self.process(stream) });
        }
    }

    pub async fn process(&self, stream: TcpStream) {
        println!("accepting connection")
    }
}
