package main

import "github.com/localcompose/locom/pkg/cmd/locom"

var (
	Name    = "locom"
	Version = "dev"
)

func main() {
	locom.Execute(Name, Version)
}
