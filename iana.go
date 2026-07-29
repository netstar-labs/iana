package iana

import (
	_ "embed"
	"sort"
	"strings"
)

// raw is the embedded IANA root-zone TLD list: a leading "# Version NNNNNNNNNN,
// Last Updated ..." comment followed by one uppercase TLD label per line. It is
// vendored and pinned in-tree (see the package doc); re-vendor deliberately, in
// lockstep with the Version() bump.
//
//go:embed data/tlds-alpha-by-domain.txt
var raw string

// set is the lowercase TLD set, built once at init. version is the pinned list
// revision parsed from the header comment. Both are immutable after init and
// safe for concurrent reads.
var (
	set     = make(map[string]struct{})
	version string
)

func init() {
	for _, line := range strings.Split(raw, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if line[0] == '#' {
			if version == "" {
				version = parseVersion(line)
			}
			continue
		}
		// TLD labels are ASCII LDH (letters, digits, hyphen); ASCII-lowercase
		// them into the set so lookups are case-insensitive and xn-- labels
		// match the A-label form idna emits.
		set[strings.ToLower(line)] = struct{}{}
	}
}

// parseVersion extracts the pinned revision from the header comment, e.g.
// "# Version 2026072702, Last Updated ..." -> "2026072702". It returns "" when
// no "Version <token>" pair is present.
func parseVersion(comment string) string {
	fields := strings.Fields(comment)
	for i, f := range fields {
		if strings.EqualFold(f, "Version") && i+1 < len(fields) {
			return strings.Trim(fields[i+1], ",")
		}
	}
	return ""
}

// IsTLD reports whether label is a current IANA root-zone top-level domain. The
// comparison is case-insensitive, so "com", "COM", and "Com" all report true, as
// does the punycode A-label form "xn--55qx5d".
//
// label must be a single DNS label. A dotted string such as "co.uk" is not a
// single TLD and reports false (only "uk" is the TLD); the empty string reports
// false. This is the root-zone existence check, not a Public Suffix List /
// registrable-domain boundary — see the package doc.
func IsTLD(label string) bool {
	_, ok := set[strings.ToLower(label)]
	return ok
}

// TLDs returns a sorted copy of every current IANA TLD, lowercased. The returned
// slice is freshly allocated on each call and owned by the caller.
func TLDs() []string {
	out := make([]string, 0, len(set))
	for tld := range set {
		out = append(out, tld)
	}
	sort.Strings(out)
	return out
}

// Count returns the number of TLDs in the pinned list.
func Count() int { return len(set) }

// Version returns the pinned IANA list revision parsed from the data file's
// header comment (e.g. "2026072702"). Stamp it wherever a TLD decision is
// persisted so a later re-vendor is detectable as skew rather than a silent flip
// (see the drift model in the package doc).
func Version() string { return version }
