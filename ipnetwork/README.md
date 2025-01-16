Fork of https://github.com/achanda/ipnetwork - commit 7e0519cd352793b17143349df0d29675b16d4fbe

ipnetwork
===
This is a library to work with IPv4 and IPv6 CIDRs in Rust

[![Build Status](https://travis-ci.org/achanda/ipnetwork.svg?branch=master)](https://travis-ci.org/achanda/ipnetwork)
[![Merit Badge](http://meritbadge.herokuapp.com/ipnetwork)](https://crates.io/crates/ipnetwork)

Run Clippy by doing
```
rustup component add clippy
cargo clippy
```

### Installation
This crate works with Cargo. Assuming you have Rust and Cargo installed, simply check out the source and run tests:
```
git clone https://github.com/achanda/ipnetwork
cd ipnetwork
cargo test
```

You can also add `ipnetwork` as a dependency to your project's `Cargo.toml`:
```toml
[dependencies]
ipnetwork = "*"
```
