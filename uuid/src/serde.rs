// Copyright 2013-2014 The Rust Project Developers.
// Copyright 2018 The Uuid Project Developers.
//
// See the COPYRIGHT file at the top-level directory of this distribution.
//
// Licensed under the Apache License, Version 2.0 <LICENSE-APACHE or
// http://www.apache.org/licenses/LICENSE-2.0> or the MIT license
// <LICENSE-MIT or http://opensource.org/licenses/MIT>, at your
// option. This file may not be copied, modified, or distributed
// except according to those terms.

use std::fmt;

use serde::{
    Deserialize, Deserializer, Serialize, Serializer,
    de::{self, Error as _},
};

use crate::{Error, Uuid};

impl Serialize for Uuid {
    fn serialize<S: Serializer>(&self, serializer: S) -> Result<S::Ok, S::Error> {
        if serializer.is_human_readable() {
            serializer.serialize_str(self.encode_lower(&mut Uuid::encode_buffer()))
        } else {
            serializer.serialize_bytes(self.as_bytes())
        }
    }
}

impl<'de> Deserialize<'de> for Uuid {
    fn deserialize<D: Deserializer<'de>>(deserializer: D) -> Result<Self, D::Error> {
        fn de_error<E: de::Error>(err: Error) -> E {
            E::custom(format_args!("UUID parsing failed: {}", err))
        }

        if deserializer.is_human_readable() {
            struct UuidVisitor;

            impl<'vi> de::Visitor<'vi> for UuidVisitor {
                type Value = Uuid;

                fn expecting(&self, formatter: &mut fmt::Formatter<'_>) -> fmt::Result {
                    write!(formatter, "a UUID string")
                }

                fn visit_str<E: de::Error>(self, value: &str) -> Result<Uuid, E> {
                    value.parse::<Uuid>().map_err(de_error)
                }

                fn visit_bytes<E: de::Error>(self, value: &[u8]) -> Result<Uuid, E> {
                    Uuid::from_slice(value).map_err(de_error)
                }

                fn visit_seq<A>(self, mut seq: A) -> Result<Uuid, A::Error>
                where
                    A: de::SeqAccess<'vi>,
                {
                    #[rustfmt::skip]
                    let bytes = [
                        match seq.next_element()? { Some(e) => e, None => return Err(A::Error::invalid_length(0, &self)) },
                        match seq.next_element()? { Some(e) => e, None => return Err(A::Error::invalid_length(1, &self)) },
                        match seq.next_element()? { Some(e) => e, None => return Err(A::Error::invalid_length(2, &self)) },
                        match seq.next_element()? { Some(e) => e, None => return Err(A::Error::invalid_length(3, &self)) },
                        match seq.next_element()? { Some(e) => e, None => return Err(A::Error::invalid_length(4, &self)) },
                        match seq.next_element()? { Some(e) => e, None => return Err(A::Error::invalid_length(5, &self)) },
                        match seq.next_element()? { Some(e) => e, None => return Err(A::Error::invalid_length(6, &self)) },
                        match seq.next_element()? { Some(e) => e, None => return Err(A::Error::invalid_length(7, &self)) },
                        match seq.next_element()? { Some(e) => e, None => return Err(A::Error::invalid_length(8, &self)) },
                        match seq.next_element()? { Some(e) => e, None => return Err(A::Error::invalid_length(9, &self)) },
                        match seq.next_element()? { Some(e) => e, None => return Err(A::Error::invalid_length(10, &self)) },
                        match seq.next_element()? { Some(e) => e, None => return Err(A::Error::invalid_length(11, &self)) },
                        match seq.next_element()? { Some(e) => e, None => return Err(A::Error::invalid_length(12, &self)) },
                        match seq.next_element()? { Some(e) => e, None => return Err(A::Error::invalid_length(13, &self)) },
                        match seq.next_element()? { Some(e) => e, None => return Err(A::Error::invalid_length(14, &self)) },
                        match seq.next_element()? { Some(e) => e, None => return Err(A::Error::invalid_length(15, &self)) },
                    ];

                    Ok(Uuid::from_bytes(bytes))
                }
            }

            deserializer.deserialize_str(UuidVisitor)
        } else {
            struct UuidBytesVisitor;

            impl<'vi> de::Visitor<'vi> for UuidBytesVisitor {
                type Value = Uuid;

                fn expecting(&self, formatter: &mut fmt::Formatter<'_>) -> fmt::Result {
                    write!(formatter, "bytes")
                }

                fn visit_bytes<E: de::Error>(self, value: &[u8]) -> Result<Uuid, E> {
                    Uuid::from_slice(value).map_err(de_error)
                }
            }

            deserializer.deserialize_bytes(UuidBytesVisitor)
        }
    }
}

// #[cfg(test)]
// mod serde_tests {
//     use super::*;

//     use serde::serde_test::{Compact, Configure, Readable, Token};

//     #[test]
//     fn test_serialize_readable_string() {
//         let uuid_str = "f9168c5e-ceb2-4faa-b6bf-329bf39fa1e4";
//         let u = Uuid::parse_str(uuid_str).unwrap();
//         serde_test::assert_tokens(&u.readable(), &[Token::Str(uuid_str)]);
//     }
// }

//     #[test]
//     fn test_deserialize_readable_compact() {
//         let uuid_bytes = b"F9168C5E-CEB2-4F";
//         let u = Uuid::from_slice(uuid_bytes).unwrap();

//         serde_test::assert_de_tokens(
//             &u.readable(),
//             &[
//                 serde_test::Token::Tuple { len: 16 },
//                 serde_test::Token::U8(uuid_bytes[0]),
//                 serde_test::Token::U8(uuid_bytes[1]),
//                 serde_test::Token::U8(uuid_bytes[2]),
//                 serde_test::Token::U8(uuid_bytes[3]),
//                 serde_test::Token::U8(uuid_bytes[4]),
//                 serde_test::Token::U8(uuid_bytes[5]),
//                 serde_test::Token::U8(uuid_bytes[6]),
//                 serde_test::Token::U8(uuid_bytes[7]),
//                 serde_test::Token::U8(uuid_bytes[8]),
//                 serde_test::Token::U8(uuid_bytes[9]),
//                 serde_test::Token::U8(uuid_bytes[10]),
//                 serde_test::Token::U8(uuid_bytes[11]),
//                 serde_test::Token::U8(uuid_bytes[12]),
//                 serde_test::Token::U8(uuid_bytes[13]),
//                 serde_test::Token::U8(uuid_bytes[14]),
//                 serde_test::Token::U8(uuid_bytes[15]),
//                 serde_test::Token::TupleEnd,
//             ],
//         );
//     }

//     #[test]
//     fn test_deserialize_readable_bytes() {
//         let uuid_bytes = b"F9168C5E-CEB2-4F";
//         let u = Uuid::from_slice(uuid_bytes).unwrap();

//         serde_test::assert_de_tokens(&u.readable(), &[serde_test::Token::Bytes(uuid_bytes)]);
//     }

//     #[test]
//     fn test_serialize_hyphenated() {
//         let uuid_str = "f9168c5e-ceb2-4faa-b6bf-329bf39fa1e4";
//         let u = Uuid::parse_str(uuid_str).unwrap();
//         serde_test::assert_ser_tokens(&u.hyphenated(), &[Token::Str(uuid_str)]);
//     }

//     #[test]
//     fn test_serialize_simple() {
//         let uuid_str = "f9168c5eceb24faab6bf329bf39fa1e4";
//         let u = Uuid::parse_str(uuid_str).unwrap();
//         serde_test::assert_ser_tokens(&u.simple(), &[Token::Str(uuid_str)]);
//     }

//     #[test]
//     fn test_serialize_urn() {
//         let uuid_str = "urn:uuid:f9168c5e-ceb2-4faa-b6bf-329bf39fa1e4";
//         let u = Uuid::parse_str(uuid_str).unwrap();
//         serde_test::assert_ser_tokens(&u.urn(), &[Token::Str(uuid_str)]);
//     }

//     #[test]
//     fn test_serialize_braced() {
//         let uuid_str = "{f9168c5e-ceb2-4faa-b6bf-329bf39fa1e4}";
//         let u = Uuid::parse_str(uuid_str).unwrap();
//         serde_test::assert_ser_tokens(&u.braced(), &[Token::Str(uuid_str)]);
//     }

//     #[test]
//     fn test_serialize_non_human_readable() {
//         let uuid_bytes = b"F9168C5E-CEB2-4F";
//         let u = Uuid::from_slice(uuid_bytes).unwrap();
//         serde_test::assert_tokens(
//             &u.compact(),
//             &[serde_test::Token::Bytes(&[
//                 70, 57, 49, 54, 56, 67, 53, 69, 45, 67, 69, 66, 50, 45, 52, 70,
//             ])],
//         );
//     }

//     #[test]
//     fn test_de_failure() {
//         serde_test::assert_de_tokens_error::<Readable<Uuid>>(
//             &[Token::Str("hello_world")],
//             "UUID parsing failed: invalid character: expected an optional prefix of `urn:uuid:` followed by [0-9a-fA-F-], found `h` at 1",
//         );

//         serde_test::assert_de_tokens_error::<Compact<Uuid>>(
//             &[Token::Bytes(b"hello_world")],
//             "UUID parsing failed: invalid length: expected 16 bytes, found 11",
//         );
//     }
// }
