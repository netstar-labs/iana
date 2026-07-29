# User guide

## Install

```
go get github.com/netstar-labs/iana
```

Build with `GOWORK=off` in a workspace that shouldn't resolve sibling modules.
There are no external dependencies to fetch.

## API

```go
func IsTLD(label string) bool   // is label a current IANA TLD? case-insensitive
func TLDs() []string            // sorted, lowercase, caller-owned copy of all TLDs
func Count() int                // number of TLDs in the pinned list
func Version() string           // the pinned list revision, e.g. "2026072702"
```

### `IsTLD(label string) bool`

Reports whether `label` is a current IANA root-zone top-level domain. Pass a
**single DNS label**, not a host.

```go
iana.IsTLD("com")         // true
iana.IsTLD("COM")         // true  — case-insensitive
iana.IsTLD("xn--55qx5d")  // true  — punycode A-label (公司)
iana.IsTLD("example")     // false — reserved name, not in the root zone
iana.IsTLD("co.uk")       // false — two labels; only "uk" is the TLD (see Scope)
iana.IsTLD("")            // false
```

Typical use is on the rightmost label of an already-parsed / canonicalized host:

```go
host := "www.example.co.uk"
labels := strings.Split(host, ".")
tld := labels[len(labels)-1]     // "uk"
if !iana.IsTLD(tld) {
    // reject: host does not end in a real TLD
}
```

### `TLDs() []string`

Returns a freshly-allocated, sorted, lowercase slice of every TLD. The slice is
yours to mutate; it does not alias the internal set.

```go
all := iana.TLDs()   // ["aaa", "aarp", ..., "zone", "zuerich", "zw"]
```

### `Count() int`

The number of TLDs in the pinned list (1438 for revision `2026072702`).

### `Version() string`

The pinned IANA list revision, parsed from the data file's header comment.

## The pin: stamp it on persisted decisions

The TLD set is a decision input. A re-vendor (below) silently changes which labels
are valid, so **whenever you persist a TLD verdict, store `Version()` alongside
it**:

```go
record.TLDValid = iana.IsTLD(tld)
record.TLDListVersion = iana.Version()   // "2026072702"
```

A later re-vendor then surfaces as detectable skew — you can see that a stored
verdict was made against an older list — rather than a silent flip. This is the
same discipline the sibling `idna` module applies to its Unicode pin.

## Re-vendoring the list (deliberate, not automatic)

The list is a committed file, never fetched at runtime. To update it:

1. Download the current list from
   `https://data.iana.org/TLD/tlds-alpha-by-domain.txt`.
2. Replace `data/tlds-alpha-by-domain.txt` with it. Keep the format: the leading
   `# Version NNNNNNNNNN, ...` header and one uppercase label per line.
3. Update the pin in `VENDOR` and the expected values in the tests
   (`TestVersionParsed` asserts the version; `TestCountMatchesFile` asserts the
   count) so the change is explicit and reviewed.
4. Re-run any persisted TLD decisions if correctness depends on the newer list —
   a changed list re-keys downstream validity verdicts.

`Version()` is single-sourced from the file's header, so step 2 moves it
automatically; steps 3–4 make the change auditable and its blast radius explicit.

## Scope: not the Public Suffix List

`iana` answers **which TLDs exist** in the IANA root zone. It does **not** compute
the eTLD+1 / registrable-domain boundary. `co.uk` is a public suffix but not a
single TLD, so `IsTLD("co.uk")` is `false`. If you need "where does the
registrable domain begin?" (that `example.co.uk` registers under `co.uk`, not
`uk`), you need a Public Suffix List — a separate data source and concern. Do not
expect PSL behavior from this package.

## Concurrency

The set is built once at package init and never mutated afterward, so all four
functions are safe for concurrent use. `TLDs()` returns a fresh copy each call.

## Gate

```
export GOWORK=off GOPRIVATE=github.com/netstar-labs
go build ./... && go vet ./... && staticcheck ./... && \
  deadcode -test ./... && test -z "$(gofmt -l .)" && go test -race ./...
```
