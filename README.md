# chrome-dump

Simple Google Chrome cookie dumper

## Usage

This project is compatible with the following target platforms:

- Windows
- Linux
- OS X

You can easily build an executable with:

```
GOOS=PLATFORM_NAME go build
```

The `dll/dll_example.go` file can help you build a DLL that can be used by other tools. To build one, run:

```
cd dll
CGO_ENABLED=1 CC=PATH_TO_CROSS_COMPILER GOOS=windows go build -buildmode c-shared -o chrome-dump.dll
```