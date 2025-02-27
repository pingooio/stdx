#[cfg(feature = "std")]
use alloc::boxed::Box;
use core::fmt::{self, Debug, Display};
#[cfg(error_generic_member_access)]
use std::error::Request;

use crate::StdError;

#[repr(transparent)]
pub struct MessageError<M>(pub M);

impl<M> Debug for MessageError<M>
where
    M: Display + Debug,
{
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        Debug::fmt(&self.0, f)
    }
}

impl<M> Display for MessageError<M>
where
    M: Display + Debug,
{
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        Display::fmt(&self.0, f)
    }
}

impl<M> StdError for MessageError<M> where M: Display + Debug + 'static {}

#[repr(transparent)]
pub struct DisplayError<M>(pub M);

impl<M> Debug for DisplayError<M>
where
    M: Display,
{
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        Display::fmt(&self.0, f)
    }
}

impl<M> Display for DisplayError<M>
where
    M: Display,
{
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        Display::fmt(&self.0, f)
    }
}

impl<M> StdError for DisplayError<M> where M: Display + 'static {}

#[cfg(feature = "std")]
#[repr(transparent)]
pub struct BoxedError(pub Box<dyn StdError + Send + Sync>);

#[cfg(feature = "std")]
impl Debug for BoxedError {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        Debug::fmt(&self.0, f)
    }
}

#[cfg(feature = "std")]
impl Display for BoxedError {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        Display::fmt(&self.0, f)
    }
}

#[cfg(feature = "std")]
impl StdError for BoxedError {
    fn source(&self) -> Option<&(dyn StdError + 'static)> {
        self.0.source()
    }

    #[cfg(error_generic_member_access)]
    fn provide<'a>(&'a self, request: &mut Request<'a>) {
        self.0.provide(request);
    }
}
