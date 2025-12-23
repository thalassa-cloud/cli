package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra/doc"
	thalassacmd "github.com/thalassa-cloud/cli/cmd"
)

// Default front matter template for generated documentation
const defaultFmTemplate = `---
linkTitle: "%s"
title: "%s"
slug: %s
url: %s
weight: %d
cascade:
  type: docs
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
		name := filepath.Base(filename)
		base := strings.TrimSuffix(name, path.Ext(name))
		displayName := strings.TrimPrefix(base, cmd.Name()+"_")

		// Create hierarchical URL structure for commands with subcommands
		var url string
		if strings.Contains(base, "_") && !strings.HasSuffix(base, "_") {
			// Check if this is a subcommand (has underscores and is not the main command)
			parts := strings.Split(base, "_")
			if len(parts) >= 3 { // tcloud_command_subcommand
				commandName := parts[1]
				subcommandPart := strings.Join(parts[2:], "_")
				url = "/docs/" + cmd.Name() + "/" + commandName + "/" + strings.ToLower(subcommandPart) + "/"
			} else {
				url = "/docs/" + cmd.Name() + "/" + strings.ToLower(base) + "/"
			}
		} else {
			url = "/docs/" + cmd.Name() + "/" + strings.ToLower(base) + "/"
		}

		weight--
		return fmt.Sprintf(defaultFmTemplate, strings.Replace(base, "_", " ", -1),
			strings.Replace(displayName, "_", " ", -1), base, url, weight)
	}

	// Generates URLs for cross-references between doc pages
	linkHandler := func(name string) string {
		base := strings.TrimSuffix(name, path.Ext(name))

		// Create hierarchical URL structure for commands with subcommands
		if strings.Contains(base, "_") && !strings.HasSuffix(base, "_") {
			// Check if this is a subcommand (has underscores and is not the main command)
			parts := strings.Split(base, "_")
			if len(parts) >= 3 { // tcloud_command_subcommand
				commandName := parts[1]
				subcommandPart := strings.Join(parts[2:], "_")
				return "/docs/" + cmd.Name() + "/" + commandName + "/" + strings.ToLower(subcommandPart) + "/"
			}
		}

		return "/docs/" + cmd.Name() + "/" + strings.ToLower(base) + "/"
	}

	// Disable auto-generated tag to avoid git conflicts
	cmd.DisableAutoGenTag = true

	// Generate documentation
	if err := doc.GenMarkdownTreeCustom(cmd, *outputDir, filePrepender, linkHandler); err != nil {
		log.Fatalf("Documentation generation failed: %v", err)
	}

	// Reorganize documentation into hierarchical structure
	if err := reorganizeDocs(*outputDir); err != nil {
		log.Fatalf("Failed to reorganize docs: %v", err)
	}

	log.Printf("Documentation successfully generated in %s", *outputDir)
}

// reorganizeDocs moves documentation files into a hierarchical structure for commands with subcommands
func reorganizeDocs(outputDir string) error {
	// Find all command files that have subcommands
	commandFiles, err := filepath.Glob(filepath.Join(outputDir, "tcloud_*.md"))
	if err != nil {
		return fmt.Errorf("failed to find command files: %w", err)
	}

	for _, file := range commandFiles {
		filename := filepath.Base(file)
		baseName := strings.TrimSuffix(filename, ".md")

		// Skip the main tcloud.md file
		if baseName == "tcloud" {
			continue
		}

		// Extract command name (e.g., "storage" from "tcloud_storage.md")
		commandName := strings.TrimPrefix(baseName, "tcloud_")

		// Find all subcommand files for this command
		pattern := filepath.Join(outputDir, "tcloud_"+commandName+"_*.md")
		subcommandFiles, err := filepath.Glob(pattern)
		if err != nil {
			return fmt.Errorf("failed to find subcommand files for %s: %w", commandName, err)
		}

		// If there are subcommands, create a directory structure
		if len(subcommandFiles) > 0 {
			commandDir := filepath.Join(outputDir, commandName)

			// Create command directory
			if err := os.MkdirAll(commandDir, 0755); err != nil {
				return fmt.Errorf("failed to create %s directory: %w", commandName, err)
			}

			// Move the main command file to _index.md in the command directory
			mainCommandFile := filepath.Join(outputDir, "tcloud_"+commandName+".md")
			if _, err := os.Stat(mainCommandFile); err == nil {
				mainIndexPath := filepath.Join(commandDir, "_index.md")
				if err := moveFile(mainCommandFile, mainIndexPath); err != nil {
					return fmt.Errorf("failed to move main command file %s to %s: %w", mainCommandFile, mainIndexPath, err)
				}
				log.Printf("Moved %s to %s", "tcloud_"+commandName+".md", mainIndexPath)
			}

			// Move subcommand files to the command directory
			for _, subFile := range subcommandFiles {
				subFilename := filepath.Base(subFile)

				// Extract subcommand name from filename (e.g., "snapshots_create" from "tcloud_storage_snapshots_create.md")
				subcommandName := strings.TrimPrefix(strings.TrimSuffix(subFilename, ".md"), "tcloud_"+commandName+"_")

				// Create subcommand directory
				subcommandDir := filepath.Join(commandDir, subcommandName)
				if err := os.MkdirAll(subcommandDir, 0755); err != nil {
					return fmt.Errorf("failed to create %s subcommand directory: %w", subcommandName, err)
				}

				// Move file as _index.md
				newPath := filepath.Join(subcommandDir, "_index.md")

				// Move the file
				if err := moveFile(subFile, newPath); err != nil {
					return fmt.Errorf("failed to move %s to %s: %w", subFile, newPath, err)
				}

				log.Printf("Moved %s to %s", subFilename, newPath)
			}
		}
	}

	return nil
}

// moveFile moves a file from src to dst
func moveFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	// Remove the original file
	return os.Remove(src)
}
