use aws_lc_rs::hmac::{self, Context, Key};

pub struct HmacSha512 {
    ctx: Context,
}

impl HmacSha512 {
    pub const SIGNATURE_SIZE: usize = 64;

    #[inline]
    pub fn sign(key: &[u8], data: &[u8]) -> [u8; 64] {
        let hmac_key = Key::new(hmac::HMAC_SHA512, key);
        return hmac::sign(&hmac_key, data).as_ref().try_into().unwrap();
    }

    #[inline]
    pub fn new(key: &[u8]) -> Self {
        let hmac_key = Key::new(hmac::HMAC_SHA512, key);
        return HmacSha512 {
            ctx: Context::with_key(&hmac_key),
        };
    }

    #[inline]
    pub fn write(&mut self, data: &[u8]) {
        self.ctx.update(data);
    }

    #[inline]
    pub fn sum(self) -> [u8; 64] {
        return self.ctx.sign().as_ref().try_into().unwrap();
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    const HELLO_WORLD_SIGNATURE: &str =
        "dce414cb1ac4e7d400ebe75f437ba90ada41c339874276b0807b7a8d9d73b56dbde7898e99c4ed92659f30ccd40c712ee517fc229012cffcd798d9ef7e357dd8";

    #[test]
    fn hello_world_signature() {
        let signature = HmacSha512::sign(b"hello world", b"hello world");

        assert_eq!(hex::encode(&signature), HELLO_WORLD_SIGNATURE);
    }

    #[test]
    fn hello_world_signer() {
        let mut hasher = HmacSha512::new(b"hello world");
        hasher.write(b"hello ");
        hasher.write(b"world");

        let signature = hasher.sum();

        assert_eq!(hex::encode(&signature), HELLO_WORLD_SIGNATURE);
    }
}
