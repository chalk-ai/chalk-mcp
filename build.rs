fn main() {
    // This is a simple build script for the manually defined protobuf types
    println!("cargo:rerun-if-changed=build.rs");
    println!("cargo:rerun-if-changed=src/proto.rs");
}