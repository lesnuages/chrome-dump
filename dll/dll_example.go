package main

import (
	"C"

	"runtime"

	"github.com/lesnuages/chrome-dump/dump"
)
import (
	"os"
)

// Build with: GOOS=windows CGO_ENABLED=1 CC=/usr/bin/x86_64-w64-mingw32-gcc go build -buildmode c-shared -o chrome-dumpe.dll

//export ChromeDump
func ChromeDump() {
	if runtime.GOOS != "windows" {
		os.Unsetenv("LD_PRELOAD")
	}
	dump.Dump()
	if runtime.GOOS != "windows" {
		os.Exit(0)
	}
}

func main() {}
