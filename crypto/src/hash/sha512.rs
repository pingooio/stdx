use aws_lc_rs::digest::{self, Context};

pub struct Sha512(Context);

impl Sha512 {
    #[inline]
    pub fn hash(output: &mut [u8], data: &[u8]) {
        let hash = digest::digest(&digest::SHA512, data);
        output.copy_from_slice(hash.as_ref());
    }

    #[inline]
    pub fn new() -> Self {
        return Sha512(Context::new(&digest::SHA512));
    }

    #[inline]
    pub fn write(&mut self, data: &[u8]) {
        self.0.update(data);
    }

    #[inline]
    pub fn sum(self, output: &mut [u8]) {
        let digest = self.0.finish();
        output.copy_from_slice(digest.as_ref());
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    const HELLO_WORLD_HASH: &str = "309ecc489c12d6eb4cc40f50c902f2b4d0ed77ee511a7c7a9bcd3ca86d4cd86f989dd35bc5ff499670da34255b45b0cfd830e81f605dcf7dc5542e93ae9cd76f";

    #[test]
    fn hello_world_hash() {
        let mut hash = [0u8; 64];
        Sha512::hash(&mut hash, b"hello world");

        assert_eq!(hex::encode(&hash), HELLO_WORLD_HASH);
    }

    #[test]
    fn hello_world_hasher() {
        let mut hasher = Sha512::new();
        hasher.write(b"hello ");
        hasher.write(b"world");

        let mut hash = [0u8; 64];
        hasher.sum(&mut hash);

        assert_eq!(hex::encode(&hash), HELLO_WORLD_HASH);
    }
}
