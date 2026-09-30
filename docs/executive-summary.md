# Executive summary

## What it is

`iana` is a small, self-contained Go library that answers one question
authoritatively: **is a given label a real top-level domain?** It embeds the
official IANA root-zone TLD list (`tlds-alpha-by-domain.txt`) directly in the
binary, pinned to a known revision, and exposes a four-function API — `IsTLD`,
`TLDs`, `Count`, `Version`. Nothing is fetched at runtime; the list ships inside
the program; there are zero external dependencies.

## Why it exists

Deciding whether a host's rightmost label is a genuine TLD is a step every URL /
host validator needs, and it is easy to get subtly wrong — by counting dots, by
using a stale or fetched-at-runtime list, or by conflating "TLD exists" with the
Public Suffix List's very different "where does the registrable domain begin?"
question. `iana` makes the TLD-validity check a single, shared, owned primitive so
that `sanitize` (host rectification, apex/TLD checks) and `normie` (URL
canonicalization) both answer it the same way, from the same pinned data.

## The one discipline that matters

The set of valid TLDs is a **decision input**, and re-vendoring it silently changes
which labels are valid. So the library exposes `Version()` — the exact revision it
is answering from — and the standing instruction is to **stamp that version wherever
a TLD verdict is persisted**. A later re-vendor then surfaces as detectable skew,
not a silent flip. This mirrors the pin/drift discipline of the sibling `idna`
module; both are governed, deliberately-updated data seams, not live feeds.

## Scope, stated plainly

- **In scope:** which top-level domains *exist* in the IANA root zone.
- **Out of scope:** the Public Suffix List / eTLD+1 registrable-domain boundary.
  `co.uk` is a public suffix but not a single TLD, so `IsTLD("co.uk")` is `false`.
  That is a different concern with a different data source.

## Where it fits

A member of a broader canonicalization family, on the **UTS-46** side alongside
a sibling host-canonicalization module: that module produces the punycode
A-label that becomes a storage key; `iana` validates the TLD that key resolves
against. UTS-39 confusable/homoglyph analysis is a separate "looks like" signal,
not part of this package.
