use sha2::Digest;

pub struct Sha256(sha2::Sha256);

#[inline]
pub fn hash_256(data: &[u8]) -> [u8; 32] {
    return sha2::Sha256::digest(data).into();
}

impl Sha256 {
    #[inline]
    pub fn new() -> Self {
        return Sha256(sha2::Sha256::new());
    }

    #[inline]
    pub fn write(&mut self, data: &[u8]) {
        self.0.update(data);
    }

    #[inline]
    pub fn sum(self) -> [u8; 32] {
        return self.0.finalize().into();
    }
}

pub struct Sha512(sha2::Sha512);

#[inline]
pub fn hash_512(data: &[u8]) -> [u8; 64] {
    return sha2::Sha512::digest(data).into();
}

impl Sha512 {
    #[inline]
    pub fn new() -> Self {
        return Sha512(sha2::Sha512::new());
    }

    #[inline]
    pub fn write(&mut self, data: &[u8]) {
        self.0.update(data);
    }

    #[inline]
    pub fn sum(self) -> [u8; 64] {
        return self.0.finalize().into();
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    const HELLO_WORLD_HASH_256: &str = "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9";
    const HELLO_WORLD_HASH_512: &str = "309ecc489c12d6eb4cc40f50c902f2b4d0ed77ee511a7c7a9bcd3ca86d4cd86f989dd35bc5ff499670da34255b45b0cfd830e81f605dcf7dc5542e93ae9cd76f";

    #[test]
    fn hello_world_hash() {
        let hash = hash_256(b"hello world");
        assert_eq!(hex::encode(&hash), HELLO_WORLD_HASH_256);

        let hash = hash_512(b"hello world");
        assert_eq!(hex::encode(&hash), HELLO_WORLD_HASH_512);
    }

    #[test]
    fn hello_world_hasher() {
        let mut hasher = Sha256::new();
        hasher.write(b"hello ");
        hasher.write(b"world");
        let hash = hasher.sum();
        assert_eq!(hex::encode(&hash), HELLO_WORLD_HASH_256);

        let mut hasher = Sha512::new();
        hasher.write(b"hello ");
        hasher.write(b"world");
        let hash = hasher.sum();
        assert_eq!(hex::encode(&hash), HELLO_WORLD_HASH_512);
    }
}
