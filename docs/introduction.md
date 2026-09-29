# iana — the list of what's real

*Every domain name on the internet ends in one of 1,438 labels. This is that
list, frozen and asked one question: is yours on it?*

Beneath every URL is a claim: that its rightmost label — `com`, `bank`, `zw`,
`xn--55qx5d` — names a top-level domain that actually exists in the root of the
DNS. Most of the time the claim is true and invisible. But attacker-controlled
input does not respect the root zone: `login.microsoft.com.verify-account.xyz` is
a real host under a real TLD (`xyz`), while `paypal.com-secure` ends in a label
(`com-secure`) that is not a TLD at all. Telling the two apart is not a matter of
counting dots — it is a matter of knowing, authoritatively, which labels the IANA
root zone actually delegates. That knowledge is a *list*, and lists drift.

`iana` is that list, owned. It embeds the official IANA root-zone TLD file —
`tlds-alpha-by-domain.txt`, one uppercase label per line under a version header —
directly into the binary with `go:embed`, pinned to a specific revision
(`2026072702`, 1,438 TLDs), and hands back a single honest answer through
`IsTLD(label)`: *is this label a current top-level domain?* Case-insensitive, so
`COM` and `com` and the punycode A-label `xn--55qx5d` all resolve the same. No
network, no file that can go missing at runtime, no third-party dependency — the
list ships inside the program, and the program cannot be surprised by an upstream
edit it never saw. `Version()` names the exact revision it is answering from, so
every verdict is traceable to a known copy of the truth.

What it deliberately is *not* is the Public Suffix List. `iana` knows that `uk` is
a TLD; it does not know — and will not pretend to know — that people register
under `co.uk` rather than directly under `uk`. `IsTLD("co.uk")` is `false`,
because `co.uk` is two labels and only one of them is a TLD. "Which labels exist
at the root" and "where does the registrable domain begin" are two different
questions with two different data sources, and conflating them is how validity
checks quietly start lying. This package answers only the first, exactly, and says
so out loud.

The operator keeps the one control that matters: the pin. The list is not fetched,
not auto-updated, not silently current — it is a committed file with a version
stamp, and it changes only when someone re-vendors it on purpose and bumps the
version to match. Because a re-vendor re-keys every downstream validity decision,
that stamp travels with the decision: store `Version()` beside a persisted verdict
and a later change becomes *detectable skew* instead of a silent flip. Freshness is
a choice made deliberately, not a surprise arriving over the wire.

`iana` is a member of a broader canonicalization family, the UTS-46 half:
where a sibling host-canonicalization module maps a host to the punycode
A-label that becomes its storage key, `iana` validates the TLD that key
resolves against. Together they are the seam `sanitize` and `normie` lean on
when they decide whether a host is real enough to canonicalize.

*Read next:* [executive summary](executive-summary.md) · [architecture](architecture.md) · [user guide](userguide.md)
