# iana

The authoritative **IANA root-zone TLD list**, embedded and pinned, behind a thin
lookup API. It answers one question — *is this label a real top-level domain?* —
for the canonicalization family's apex/validity boundary. **Zero external
dependencies.**

```go
iana.IsTLD("com")          // true
iana.IsTLD("XN--55QX5D")   // true  — case-insensitive; punycode A-label (公司)
iana.IsTLD("example")      // false — reserved name, not in the root zone
iana.IsTLD("co.uk")        // false — a public suffix, but NOT a single TLD
iana.IsTLD("")             // false

iana.Count()               // 1438
iana.Version()             // "2026072702"  — the pin; stamp it on stored decisions
iana.TLDs()                // ["aaa", "aarp", ..., "zw"]  — sorted, lowercase copy
```

## Data flow

```
   host's rightmost label ──▶ iana.IsTLD ──▶ is this a real TLD? (yes / no)
                                  │
   data/tlds-alpha-by-domain.txt ─┘  (go:embed, pinned to 2026072702, 1438 TLDs)
                                     Version() stamps the revision on every
                                     persisted TLD decision → re-vendor is detectable
```

## Scope: the root zone, not the Public Suffix List

This is the IANA **root-zone TLD list** — which TLDs *exist* (`com`, `uk`,
`xn--55qx5d`). It is **not** the Public Suffix List and does not compute the
eTLD+1 / registrable-domain boundary. `co.uk` is a public suffix but not a single
TLD, so `IsTLD("co.uk")` is `false`. If you need "where does the registrable
domain begin?", that is a separate concern with a separate data source. See
[docs/architecture.md](docs/architecture.md#scope-boundary).

## The pin

The TLD set is a decision input, and a re-vendor silently changes which labels are
valid. `Version()` returns the pinned list revision parsed from the data file's
header; stamp it wherever a TLD verdict is persisted so a later re-vendor shows up
as skew rather than a silent flip. This mirrors the drift discipline of the
sibling [`idna`](https://github.com/netstar-labs/idna) module.

## Where this sits

A member of the **canonicalization** family
([handbook/families/canon.md](https://github.com/netstar-labs/handbook/blob/main/families/canon.md)),
feeding the **UTS-46** side: `idna` maps a host to its A-label (the key), `iana`
validates the TLD that key resolves against. Both `sanitize` (host rectify +
TLD/apex) and `normie` (URL canonicalization) consume the pair. UTS-39 confusable
analysis (`homoglyph`) is the separate "what does it look like?" signal.

## Docs

- **Start here** — [introduction](docs/introduction.md) · [executive summary](docs/executive-summary.md)
- **Deep dive** — [architecture](docs/architecture.md)
- **Operations** — [user guide](docs/userguide.md)
- **Examples** — [example/README.md](example/README.md)

## Layout

| Path | Purpose |
|---|---|
| [iana.go](iana.go) | the library — `IsTLD`, `TLDs`, `Count`, `Version`; embeds and parses the list at init |
| [doc.go](doc.go) | package doc — the scope boundary (root zone, not PSL) and the pin/stamp discipline |
| [data/tlds-alpha-by-domain.txt](data/tlds-alpha-by-domain.txt) | the pinned IANA source, embedded via `go:embed` (not testdata — it ships) |
| [VENDOR](VENDOR) | the pin record + re-vendor discipline |

Go module `github.com/netstar-labs/iana`. Standard library only. Build with
`GOWORK=off`.
