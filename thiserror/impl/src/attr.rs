use std::collections::BTreeSet as Set;

use proc_macro2::{Delimiter, Group, Span, TokenStream, TokenTree};
use quote::{ToTokens, format_ident, quote};
use syn::{
    Attribute, Error, Ident, Index, LitInt, LitStr, Meta, Result, Token, braced, bracketed, parenthesized,
    parse::{ParseStream, discouraged::Speculative},
    token,
};

pub struct Attrs<'a> {
    pub display: Option<Display<'a>>,
    pub source: Option<&'a Attribute>,
    pub backtrace: Option<&'a Attribute>,
    pub from: Option<&'a Attribute>,
    pub transparent: Option<Transparent<'a>>,
}

#[derive(Clone)]
pub struct Display<'a> {
    pub original: &'a Attribute,
    pub fmt: LitStr,
    pub args: TokenStream,
    pub requires_fmt_machinery: bool,
    pub has_bonus_display: bool,
    pub implied_bounds: Set<(usize, Trait)>,
}

#[derive(Copy, Clone)]
pub struct Transparent<'a> {
    pub original: &'a Attribute,
    pub span: Span,
}

#[derive(Copy, Clone, Eq, PartialEq, Ord, PartialOrd, Debug)]
pub enum Trait {
    Debug,
    Display,
    Octal,
    LowerHex,
    UpperHex,
    Pointer,
    Binary,
    LowerExp,
    UpperExp,
}

pub fn get(input: &[Attribute]) -> Result<Attrs> {
    let mut attrs = Attrs {
        display: None,
        source: None,
        backtrace: None,
        from: None,
        transparent: None,
    };

    for attr in input {
        if attr.path().is_ident("error") {
            parse_error_attribute(&mut attrs, attr)?;
        } else if attr.path().is_ident("source") {
            attr.meta.require_path_only()?;
            if attrs.source.is_some() {
                return Err(Error::new_spanned(attr, "duplicate #[source] attribute"));
            }
            attrs.source = Some(attr);
        } else if attr.path().is_ident("backtrace") {
            attr.meta.require_path_only()?;
            if attrs.backtrace.is_some() {
                return Err(Error::new_spanned(attr, "duplicate #[backtrace] attribute"));
            }
            attrs.backtrace = Some(attr);
        } else if attr.path().is_ident("from") {
            match attr.meta {
                Meta::Path(_) => {}
                Meta::List(_) | Meta::NameValue(_) => {
                    // Assume this is meant for derive_more crate or something.
                    continue;
                }
            }
            if attrs.from.is_some() {
                return Err(Error::new_spanned(attr, "duplicate #[from] attribute"));
            }
            attrs.from = Some(attr);
        }
    }

    Ok(attrs)
}

fn parse_error_attribute<'a>(attrs: &mut Attrs<'a>, attr: &'a Attribute) -> Result<()> {
    syn::custom_keyword!(transparent);

    attr.parse_args_with(|input: ParseStream| {
        if let Some(kw) = input.parse::<Option<transparent>>()? {
            if attrs.transparent.is_some() {
                return Err(Error::new_spanned(attr, "duplicate #[error(transparent)] attribute"));
            }
            attrs.transparent = Some(Transparent {
                original: attr,
                span: kw.span,
            });
            return Ok(());
        }

        let fmt: LitStr = input.parse()?;

        let ahead = input.fork();
        ahead.parse::<Option<Token![,]>>()?;
        let args = if ahead.is_empty() {
            input.advance_to(&ahead);
            TokenStream::new()
        } else {
            parse_token_expr(input, false)?
        };

        let requires_fmt_machinery = !args.is_empty();

        let display = Display {
            original: attr,
            fmt,
            args,
            requires_fmt_machinery,
            has_bonus_display: false,
            implied_bounds: Set::new(),
        };
        if attrs.display.is_some() {
            return Err(Error::new_spanned(attr, "only one #[error(...)] attribute is allowed"));
        }
        attrs.display = Some(display);
        Ok(())
    })
}

fn parse_token_expr(input: ParseStream, mut begin_expr: bool) -> Result<TokenStream> {
    let mut tokens = Vec::new();
    while !input.is_empty() {
        if begin_expr && input.peek(Token![.]) {
            if input.peek2(Ident) {
                input.parse::<Token![.]>()?;
                begin_expr = false;
                continue;
            }
            if input.peek2(LitInt) {
                input.parse::<Token![.]>()?;
                let int: Index = input.parse()?;
                let ident = format_ident!("_{}", int.index, span = int.span);
                tokens.push(TokenTree::Ident(ident));
                begin_expr = false;
                continue;
            }
        }

        begin_expr = input.peek(Token![break])
            || input.peek(Token![continue])
            || input.peek(Token![if])
            || input.peek(Token![in])
            || input.peek(Token![match])
            || input.peek(Token![mut])
            || input.peek(Token![return])
            || input.peek(Token![while])
            || input.peek(Token![+])
            || input.peek(Token![&])
            || input.peek(Token![!])
            || input.peek(Token![^])
            || input.peek(Token![,])
            || input.peek(Token![/])
            || input.peek(Token![=])
            || input.peek(Token![>])
            || input.peek(Token![<])
            || input.peek(Token![|])
            || input.peek(Token![%])
            || input.peek(Token![;])
            || input.peek(Token![*])
            || input.peek(Token![-]);

        let token: TokenTree = if input.peek(token::Paren) {
            let content;
            let delimiter = parenthesized!(content in input);
            let nested = parse_token_expr(&content, true)?;
            let mut group = Group::new(Delimiter::Parenthesis, nested);
            group.set_span(delimiter.span.join());
            TokenTree::Group(group)
        } else if input.peek(token::Brace) {
            let content;
            let delimiter = braced!(content in input);
            let nested = parse_token_expr(&content, true)?;
            let mut group = Group::new(Delimiter::Brace, nested);
            group.set_span(delimiter.span.join());
            TokenTree::Group(group)
        } else if input.peek(token::Bracket) {
            let content;
            let delimiter = bracketed!(content in input);
            let nested = parse_token_expr(&content, true)?;
            let mut group = Group::new(Delimiter::Bracket, nested);
            group.set_span(delimiter.span.join());
            TokenTree::Group(group)
        } else {
            input.parse()?
        };
        tokens.push(token);
    }
    Ok(TokenStream::from_iter(tokens))
}

impl ToTokens for Display<'_> {
    fn to_tokens(&self, tokens: &mut TokenStream) {
        let fmt = &self.fmt;
        let args = &self.args;

        // Currently `write!(f, "text")` produces less efficient code than
        // `f.write_str("text")`. We recognize the case when the format string
        // has no braces and no interpolated values, and generate simpler code.
        tokens.extend(if self.requires_fmt_machinery {
            quote! {
                ::core::write!(__formatter, #fmt #args)
            }
        } else {
            quote! {
                __formatter.write_str(#fmt)
            }
        });
    }
}

impl ToTokens for Trait {
    fn to_tokens(&self, tokens: &mut TokenStream) {
        let trait_name = match self {
            Trait::Debug => "Debug",
            Trait::Display => "Display",
            Trait::Octal => "Octal",
            Trait::LowerHex => "LowerHex",
            Trait::UpperHex => "UpperHex",
            Trait::Pointer => "Pointer",
            Trait::Binary => "Binary",
            Trait::LowerExp => "LowerExp",
            Trait::UpperExp => "UpperExp",
        };
        let ident = Ident::new(trait_name, Span::call_site());
        tokens.extend(quote!(::core::fmt::#ident));
    }
}
