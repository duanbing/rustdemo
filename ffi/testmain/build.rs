fn main() {
    let path = "../golib";
    let lib = "math";

    println!("cargo:rustc-link-search=native={}", path);
    println!("cargo:rustc-link-lib=static={}", lib);
}
