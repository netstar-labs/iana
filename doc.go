// Package iana is the UTS-46 TLD-validity seam for the netstar canonicalization
// family: it answers "is this label a real top-level domain?" against the
// authoritative IANA root-zone TLD list, embedded and pinned in-tree, behind a
// thin owned API with zero external dependencies.
//
// The list ([data/tlds-alpha-by-domain.txt], one uppercase TLD per line under a
// "# Version NNNNNNNNNN" header) is vendored and embedded with go:embed, so the
// set of valid TLDs is frozen against upstream drift and the module builds with
// nothing but the standard library. It is pinned to a specific IANA revision;
// [Version] returns that revision.
//
// # Scope: the root zone, not the Public Suffix List
//
// This package answers exactly one question — does a top-level domain of this
// name EXIST in the IANA root zone (is "com", "org", "uk", "xn--55qx5d" a real
// TLD)? That is the apex/validity boundary sanitize and normie need when they
// decide whether a host's rightmost label is a genuine TLD.
//
// It is NOT the Public Suffix List and does NOT compute the eTLD+1 /
// registrable-domain boundary. "co.uk" is a public suffix but it is not a single
// TLD, so [IsTLD]("co.uk") is false — "co.uk" is two labels, and only "uk" is the
// TLD here. Callers that need "where does the registrable domain begin?" (that
// "example.co.uk" registers under "co.uk", not "uk") want a PSL, which is a
// separate concern and a separate data source. Do not expect PSL behavior here.
//
// # The pin and the stamp
//
// Mirroring the drift discipline of the sibling idna module: the TLD set is a
// decision input, and a re-vendor silently changes which labels are valid. Stamp
// [Version] wherever a TLD decision is persisted, so a later re-vendor is
// detectable as skew (a stored "valid TLD" verdict was made against a known list
// revision) rather than a silent flip. The set is immutable after init and safe
// for concurrent use.
//
// This is a member of the canonicalization family (see
// netstar-labs/handbook/families/canon.md), feeding the UTS-46 side alongside
// idna: idna maps a host to its A-label (the key), iana validates the TLD that
// key resolves against. UTS-39 confusable analysis (homoglyph) is the separate
// "what does it look like?" signal and lives elsewhere.
package iana
