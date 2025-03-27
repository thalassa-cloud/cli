package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra/doc"
	thalassacmd "github.com/thalassa-cloud/cli/cmd"
)

// Default front matter template for generated documentation
const defaultFmTemplate = `---
date: %s
linkTitle: "%s"
title: "%s"
slug: %s
url: %s
weight: %d
---
`

func main() {
	// Command line flags for customization
	outputDir := flag.String("output", "./docs/tcloud", "Output directory for generated docs")
	startWeight := flag.Int("weight", 10000, "Starting weight for documentation entries")
	flag.Parse()

	// Ensure output directory exists
	if err := os.MkdirAll(*outputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	cmd := thalassacmd.RootCmd
	weight := *startWeight

	// Prepends front matter to each generated markdown file
	filePrepender := func(filename string) string {
		now := time.Now().Format(time.RFC3339)
		name := filepath.Base(filename)
		base := strings.TrimSuffix(name, path.Ext(name))
		displayName := strings.TrimPrefix(base, cmd.Name()+"_")
		url := "/docs/" + cmd.Name() + "/" + strings.ToLower(base) + "/"
		weight--
		return fmt.Sprintf(defaultFmTemplate, now, strings.Replace(base, "_", " ", -1),
			strings.Replace(displayName, "_", " ", -1), base, url, weight)
	}

	// Generates URLs for cross-references between doc pages
	linkHandler := func(name string) string {
		base := strings.TrimSuffix(name, path.Ext(name))
		return "/docs/" + cmd.Name() + "/" + strings.ToLower(base) + "/"
	}

	// Disable auto-generated tag to avoid git conflicts
	cmd.DisableAutoGenTag = true

	// Generate documentation
	if err := doc.GenMarkdownTreeCustom(cmd, *outputDir, filePrepender, linkHandler); err != nil {
		log.Fatalf("Documentation generation failed: %v", err)
	}

	log.Printf("Documentation successfully generated in %s", *outputDir)
}
