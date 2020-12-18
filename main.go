package main

import (
	"flag"

	"github.com/lesnuages/chrome-dump/dump"
)

func main() {
	var remote string
	flag.StringVar(&remote, "remote", "", "WS url")
	flag.Parse()
	// dump.Dump(remote)
	dump.Spy(remote)
}
