package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/localcompose/locom/pkg/cmd/locom"

	"github.com/spf13/cobra/doc"
)

func main() {
	rootCmd := locom.NewRootCmd()

	// Generate Markdown docs
	err := doc.GenMarkdownTree(rootCmd, "./docs")
	if err != nil {
		log.Fatal(err)
	}
}

func main1() {
	rootCmd := locom.NewRootCmd() // however you construct your root

	// 1. Generate full docs into ./docs
	outDir := filepath.Join(".", "docs")
	os.MkdirAll(outDir, 0755)
	if err := doc.GenMarkdownTree(rootCmd, outDir); err != nil {
		log.Fatalf("gen docs: %v", err)
	}

	// 2. Generate a short overview for README.md
	var sb strings.Builder
	doc.GenMarkdown(rootCmd, &sb)

	// Insert between markers
	readme := "README.md"
	input, err := os.ReadFile(readme)
	if err != nil {
		log.Fatalf("read readme: %v", err)
	}

	lines := strings.Split(string(input), "\n")
	var out []string
	inside := false
	for _, line := range lines {
		if strings.Contains(line, "<!-- commands:start -->") {
			inside = true
			out = append(out, line, sb.String())
			continue
		}
		if strings.Contains(line, "<!-- commands:end -->") {
			inside = false
		}
		if !inside {
			out = append(out, line)
		}
	}
	os.WriteFile(readme, []byte(strings.Join(out, "\n")), 0644)

	fmt.Println("Docs updated.")
}
