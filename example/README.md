# Examples

Runnable programs showing how to use `github.com/netstar-labs/iana`.

| Example | What it shows | Run |
|---|---|---|
| [istld](istld/main.go) | `IsTLD` per label; `Version` / `Count` / `TLDs` with no args | `go run ./example/istld com org co.uk example xn--55qx5d` |

Expected output for the command above:

```
com              TLD
org              TLD
co.uk            not a TLD
example          not a TLD
xn--55qx5d       TLD
```

With no arguments it prints the pin instead:

```
go run ./example/istld
```

```
IANA root-zone TLD list: version 2026072702, 1438 TLDs
first: [aaa aarp abb] ... last: [zone zuerich zw]
```
