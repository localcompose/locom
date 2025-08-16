package main

import "github.com/localcompose/locom/cmd"

var (
	Name    = "locom"
	Version = "dev"
)

func main() {
	cmd.Execute(Name, Version)
}
