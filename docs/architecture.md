# Architecture

`iana` is a flat, single-package library. It has exactly one moving part — a set
built once from an embedded file — and a four-function surface over it.

## Subsystems

```
   data/tlds-alpha-by-domain.txt        the pinned IANA source (committed, ships)
            │  //go:embed  →  var raw string
            ▼
   init(): parse                        split on \n; skip blank lines; the "#"
            │                           header yields Version(); each remaining
            │                           line is ASCII-lowercased into the set
            ▼
   set  map[string]struct{}   +   version string      (immutable after init)
            │
            ├── IsTLD(label)  →  set[ToLower(label)] present?      O(1)
            ├── TLDs()        →  sorted lowercase copy             O(n log n)
            ├── Count()       →  len(set)                          O(1)
            └── Version()     →  version                           O(1)
```

### The embedded source

`data/tlds-alpha-by-domain.txt` is the official IANA list: a leading
`# Version NNNNNNNNNN, Last Updated ...` comment followed by one **uppercase** TLD
label per line, including punycode `XN--…` labels for internationalized TLDs. It
is embedded with `//go:embed` into a package-level `string`, so it is part of the
compiled artifact — there is no runtime file, no network fetch, and no way for the
list to differ between build and run.

It lives in `data/`, not `testdata/`, on purpose: it is the shipped source of
truth, not a test fixture. `.gitignore` ignores `testdata/` but keeps `data/`
committed.

### Parse at init

A single `init()` walks the embedded text once:

- Blank lines (after trimming) are skipped.
- The first `#` line is parsed for the version token (`Version 2026072702` →
  `2026072702`); later `#` lines, if any, are ignored.
- Every other line is a TLD label, ASCII-lowercased into `set` (a
  `map[string]struct{}`). Labels are LDH (letters, digits, hyphen), so
  ASCII-lowercasing is exact and makes both lookups and stored punycode
  A-labels case-insensitive.

After `init` the `set` and `version` are never written again, so all reads are
safe for concurrent use without synchronization.

### The lookup surface

- `IsTLD(label)` — one `map` probe on the lowercased label. A dotted string like
  `"co.uk"` is simply not a key, so it returns `false` for free — no dot-splitting
  logic, because that is a different (PSL) question.
- `TLDs()` — allocates a fresh slice, copies the keys, sorts, returns it. The copy
  is caller-owned; mutating it cannot corrupt the set.
- `Count()` / `Version()` — direct reads of `len(set)` and the parsed pin.

## Design trade-offs

- **Set over sorted-slice-plus-binary-search.** The whole point is membership;
  a map gives O(1) lookups and the simplest possible `IsTLD`. `TLDs()` pays the
  sort only when someone actually wants the enumeration.
- **`init()` over lazy `sync.Once`.** The data is ~9 KB and always needed by any
  caller of `IsTLD`; parsing eagerly at package load is simpler than lazy
  initialization and removes a branch from every lookup. There is no
  configuration to defer.
- **Parse the version from the file, don't hard-code it.** `Version()` reads the
  header comment, so the pin is single-sourced in the data file itself; a
  re-vendor that replaces the file also moves the version, and the test asserts
  the expected value so a mismatch is caught.
- **Zero dependencies, stdlib only.** `embed`, `sort`, `strings`. Nothing else.

## Scope boundary

This package is the IANA **root-zone TLD list** — the set of TLDs that *exist*. It
is deliberately **not** the Public Suffix List (PSL) and does not compute the
eTLD+1 / registrable-domain boundary:

| Question | Answered here? | Where it belongs |
|---|---|---|
| Does the TLD `uk` exist? | Yes — `IsTLD("uk")` → true | this package |
| Is `co.uk` a single TLD? | Answered `false` — it is two labels | this package (correctly negative) |
| Does one register under `co.uk` rather than `uk`? | **No** | a Public Suffix List, separate data + module |

Conflating the two is a classic source of validity bugs, so the boundary is stated
in the package doc, the README, and here. `iana` answers root-zone existence,
nothing more.

## Family context

`iana` is a member of the **canonicalization** family (see
[handbook/families/canon.md](https://github.com/netstar-labs/handbook/blob/main/families/canon.md)),
on the **UTS-46** side:

```
  host ─▶ idna.ToASCII ─▶ A-label (the storage key) ─┐
                                                      ├─ UTS-46: what does the URL resolve to?
  A-label's rightmost label ─▶ iana.IsTLD ─▶ real TLD?┘
```

`idna` (pinned to a Unicode version) produces the A-label that becomes a storage
key; `iana` (pinned to an IANA root-zone revision) validates the TLD that key
resolves against. Both are **pinned, deliberately-updated data seams** sharing the
same drift discipline: a change re-keys downstream decisions, so it is stamped and
migrated, never auto-updated. UTS-39 confusable/homoglyph analysis is the separate
"what does it look like?" signal and does not live here.
