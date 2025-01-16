use aws_lc_rs::digest::{self, Context};

#[derive(Clone)]
pub struct Sha256 {
    ctx: Context,
}

impl Sha256 {
    pub const HASH_SIZE: usize = 32;

    #[inline]
    pub fn hash(data: &[u8]) -> [u8; 32] {
        return digest::digest(&digest::SHA256, data).as_ref().try_into().unwrap();
    }

    #[inline]
    pub fn new() -> Self {
        return Sha256 {
            ctx: Context::new(&digest::SHA256),
        };
    }

    #[inline]
    pub fn write(&mut self, data: &[u8]) {
        self.ctx.update(data);
    }

    #[inline]
    pub fn sum(self) -> [u8; 32] {
        return self.ctx.finish().as_ref().try_into().unwrap();
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    const HELLO_WORLD_HASH: &str = "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9";

    #[test]
    fn hello_world_hash() {
        let hash = Sha256::hash(b"hello world");

        assert_eq!(hex::encode(&hash), HELLO_WORLD_HASH);
    }

    #[test]
    fn hello_world_hasher() {
        let mut hasher = Sha256::new();
        hasher.write(b"hello ");
        hasher.write(b"world");

        let hash = hasher.sum();

        assert_eq!(hex::encode(&hash), HELLO_WORLD_HASH);
    }
}
