// Command istld reports, for each argument, whether it is a current IANA
// top-level domain. With no arguments it prints the pinned list's version and
// count and a small sample of TLDs.
//
//	go run ./example/istld com org co.uk example xn--55qx5d
package main

import (
	"fmt"
	"os"

	"github.com/netstar-labs/iana"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		all := iana.TLDs()
		fmt.Printf("IANA root-zone TLD list: version %s, %d TLDs\n", iana.Version(), iana.Count())
		fmt.Printf("first: %v ... last: %v\n", all[:3], all[len(all)-3:])
		fmt.Println("\npass labels as arguments to test them, e.g.:")
		fmt.Println("  go run ./example/istld com org co.uk example xn--55qx5d")
		return
	}

	for _, label := range args {
		mark := "not a TLD"
		if iana.IsTLD(label) {
			mark = "TLD"
		}
		fmt.Printf("%-16s %s\n", label, mark)
	}
}
