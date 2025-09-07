package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/localcompose/locom/pkg/cmd/locom"
)

var (
	Name    = "locom"
	Version = "dev"

	supportedOS = []string{"linux", "darwin", "windows"}
)

func main() {
	if !isSupported(runtime.GOOS) {
		fmt.Fprintf(os.Stderr,
			"Sorry, your operating system (%s) is not supported.\nSupported systems are: %v\n",
			runtime.GOOS, supportedOS)
		os.Exit(1)
	}

	locom.Execute(Name, Version)
}

func isSupported(goos string) bool {
	for _, s := range supportedOS {
		if goos == s {
			return true
		}
	}
	return false
}
