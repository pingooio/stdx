use aws_lc_rs::constant_time;

pub mod aead;
pub mod hash;
pub mod hmac;

#[inline]
pub fn constant_time_compare(a: &[u8], b: &[u8]) -> bool {
    return constant_time::verify_slices_are_equal(a, b).is_ok();
}
