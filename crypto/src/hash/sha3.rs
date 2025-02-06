use aws_lc_rs::digest::{self, Context};

pub struct Sha3_256(Context);

#[inline]
pub fn hash_256(data: &[u8]) -> [u8; 32] {
    let mut ret = [0u8; 32];
    let digest = digest::digest(&digest::SHA3_256, data);
    // TODO: is there a better way to transform a digest into an array?
    ret.copy_from_slice(digest.as_ref());
    return ret;
}

impl Sha3_256 {
    #[inline]
    pub fn new() -> Self {
        return Sha3_256(Context::new(&digest::SHA3_256));
    }

    #[inline]
    pub fn write(&mut self, data: &[u8]) {
        self.0.update(data);
    }

    #[inline]
    pub fn sum(self) -> [u8; 32] {
        let mut ret = [0u8; 32];
        let digest = self.0.finish();
        // TODO: is there a better way to transform a digest into an array?
        ret.copy_from_slice(digest.as_ref());
        return ret;
    }
}

pub struct Sha3_512(Context);

#[inline]
pub fn hash_512(data: &[u8]) -> [u8; 64] {
    let mut ret = [0u8; 64];
    let digest = digest::digest(&digest::SHA3_512, data);
    // TODO: is there a better way to transform a digest into an array?
    ret.copy_from_slice(digest.as_ref());
    return ret;
}

impl Sha3_512 {
    #[inline]
    pub fn new() -> Self {
        return Sha3_512(Context::new(&digest::SHA3_512));
    }

    #[inline]
    pub fn write(&mut self, data: &[u8]) {
        self.0.update(data);
    }

    #[inline]
    pub fn sum(self) -> [u8; 64] {
        let mut ret = [0u8; 64];
        let digest = self.0.finish();
        // TODO: is there a better way to transform a digest into an array?
        ret.copy_from_slice(digest.as_ref());
        return ret;
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    const HELLO_WORLD_HASH_256: &str = "644bcc7e564373040999aac89e7622f3ca71fba1d972fd94a31c3bfbf24e3938";
    const HELLO_WORLD_HASH_512: &str = "840006653e9ac9e95117a15c915caab81662918e925de9e004f774ff82d7079a40d4d27b1b372657c61d46d470304c88c788b3a4527ad074d1dccbee5dbaa99a";

    #[test]
    fn hello_world_hash() {
        let hash = hash_256(b"hello world");
        assert_eq!(hex::encode(&hash), HELLO_WORLD_HASH_256);

        let hash = hash_512(b"hello world");
        assert_eq!(hex::encode(&hash), HELLO_WORLD_HASH_512);
    }

    #[test]
    fn hello_world_hasher() {
        let mut hasher = Sha3_256::new();
        hasher.write(b"hello ");
        hasher.write(b"world");
        let hash = hasher.sum();
        assert_eq!(hex::encode(&hash), HELLO_WORLD_HASH_256);

        let mut hasher = Sha3_512::new();
        hasher.write(b"hello ");
        hasher.write(b"world");
        let hash = hasher.sum();
        assert_eq!(hex::encode(&hash), HELLO_WORLD_HASH_512);
    }
}
