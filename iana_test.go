package iana

import (
	"sort"
	"strings"
	"testing"
)

func TestIsTLD(t *testing.T) {
	cases := []struct {
		label string
		want  bool
	}{
		{"com", true},
		{"org", true},
		{"uk", true},
		{"xn--55qx5d", true}, // punycode A-label (公司), stored uppercase in the list

		// case-insensitivity: same three labels in mixed / upper case
		{"COM", true},
		{"Org", true},
		{"XN--55QX5D", true},

		// not TLDs
		{"example", false}, // reserved special-use name, not in the root zone
		{"", false},        // empty is not a label
		{"co.uk", false},   // a public suffix, but NOT a single TLD (two labels)
		{"example.com", false},
		{"notarealtld", false},
	}
	for _, c := range cases {
		if got := IsTLD(c.label); got != c.want {
			t.Errorf("IsTLD(%q) = %v, want %v", c.label, got, c.want)
		}
	}
}

func TestTLDsSortedLowercaseAndComplete(t *testing.T) {
	tlds := TLDs()

	if len(tlds) != Count() {
		t.Fatalf("len(TLDs()) = %d, want Count() = %d", len(tlds), Count())
	}
	if !sort.StringsAreSorted(tlds) {
		t.Error("TLDs() is not sorted")
	}
	for _, tld := range tlds {
		if tld != strings.ToLower(tld) {
			t.Errorf("TLDs() returned non-lowercase entry %q", tld)
		}
		if !IsTLD(tld) {
			t.Errorf("TLDs() entry %q is not reported by IsTLD", tld)
		}
	}

	// TLDs returns a fresh, caller-owned copy: mutating it must not affect the set.
	if len(tlds) > 0 {
		tlds[0] = "MUTATED"
		if !sort.StringsAreSorted(TLDs()) {
			t.Error("mutating the returned slice affected a later TLDs() call")
		}
	}
}

func TestCountMatchesFile(t *testing.T) {
	// Independently recount the non-comment, non-blank lines in the embedded
	// file and confirm Count() agrees.
	var n int
	for _, line := range strings.Split(raw, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		n++
	}
	if Count() != n {
		t.Errorf("Count() = %d, want %d (non-comment lines in the embedded file)", Count(), n)
	}
	if Count() != 1438 {
		t.Errorf("Count() = %d, want 1438 (the pinned %s list)", Count(), Version())
	}
}

func TestVersionParsed(t *testing.T) {
	if got, want := Version(), "2026072702"; got != want {
		t.Errorf("Version() = %q, want %q (bump in lockstep with a re-vendor)", got, want)
	}
}

func TestParseVersion(t *testing.T) {
	cases := []struct {
		comment string
		want    string
	}{
		{"# Version 2026072702, Last Updated Tue Jul 28 07:07:01 2026 UTC", "2026072702"},
		{"# version 123", "123"}, // case-insensitive keyword match
		{"# just a comment", ""}, // no "version" token -> ""
		{"# Version", ""},        // keyword with no following token -> ""
		{"#", ""},
	}
	for _, c := range cases {
		if got := parseVersion(c.comment); got != c.want {
			t.Errorf("parseVersion(%q) = %q, want %q", c.comment, got, c.want)
		}
	}
}
