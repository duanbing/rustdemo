use std::os::unix::net::{UnixStream};
use std::io::{Write,Read};
fn main() -> std::io::Result<()> {
    let mut socket = UnixStream::connect("/tmp/my.sock")?;
    let msg = String::from("nihao!!!!!cac acamcakc app 11234");

    let siz = std::format!("{:0>32}",msg.len());
    println!("size : {}", siz);

    socket.write_all(siz.as_bytes());
    socket.write_all(msg.as_bytes());


    let mut resp = String::new();
    socket.read_to_string(&mut resp);
    println!("{}", resp);
    Ok(())
}

