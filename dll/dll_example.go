package main

import (
	"C"

	"runtime"

	"github.com/lesnuages/chrome-dump/dump"
)
import (
	"os"
	"strings"
)

// Build with: GOOS=windows CGO_ENABLED=1 CC=/usr/bin/x86_64-w64-mingw32-gcc go build -buildmode c-shared -o chrome-dumpe.dll

func getRemote(p string) string {
	var remote string
	s := strings.Split(p, "-remote ")
	if len(s) > 1 {
		remote = s[1]
	}
	return remote
}

//export ChromeDump
func ChromeDump() {
	params := os.Getenv("LD_PARAMS")
	if runtime.GOOS != "windows" {
		os.Unsetenv("LD_PRELOAD")
	}
	dump.Dump(getRemote(params))
	if runtime.GOOS != "windows" {
		os.Exit(0)
	}
}

func main() {}
