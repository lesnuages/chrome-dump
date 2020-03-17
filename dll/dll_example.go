package main

import (
	"C"

	"github.com/lesnuages/chrome-dump/dump"
)
import (
	"os"
)

// Build with: GOOS=windows CGO_ENABLED=1 CC=/usr/bin/x86_64-w64-mingw32-gcc go build -buildmode c-shared -o chrome-dumpe.dll

//export ChromeDump
func ChromeDump() {
	os.Unsetenv("LD_PRELOAD")
	dump.Dump()
	os.Exit(0)
}

func main() {}
