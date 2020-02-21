use std::thread;
use std::os::unix::net::{UnixStream, UnixListener};
use std::io::{Read};

fn handle_client(mut stream: UnixStream) {
    println!("handle client ... ");
    // 读取长度
    let mut sz = 0;
    {
        let mut buf: [u8; 32] = [0; 32];
        let mut handle = stream.by_ref().take(32);
        handle.read(&mut buf);
        let sz_str = String::from_utf8(buf.to_vec()).unwrap();
        println!("{:?}",sz_str);
        sz = sz_str.parse::<u32>().unwrap();
        println!("read length: {}", sz);
    }

    let mut handle = stream.take(sz as u64);
    let mut buf = vec![0u8; sz as usize];
    handle.read(&mut buf);
    println!("read content: {:?}", String::from_utf8(buf.to_vec()));

    //读取实际字符串


}

fn main() -> std::io::Result<()> {
    let path = "/tmp/my.sock";

    if std::path::Path::new(path).exists() {
       std::fs::remove_dir(path);
    }
    let listener = UnixListener::bind(path)?;
    for stream in listener.incoming() {
        match stream {
            Ok(stream) => {
                thread::spawn(||handle_client(stream));
            },
            Err(e) => {
                break;
            },
        }
    }
    println!("Hello, world!");
    Ok(())
}
