package main

import (
	"fmt"
	"os"

	"github.com/customrealms/cli/internal/build"
	"github.com/customrealms/cli/internal/project"
	"gopkg.in/yaml.v3"
)

type YmlCmd struct {
	ProjectDir string `name:"project" short:"p" usage:"plugin project directory" optional:""`
}

func (c *YmlCmd) Run() error {
	// Root context for the CLI
	_, cancel := rootContext()
	defer cancel()

	// Default to the current working directory
	if c.ProjectDir == "" {
		c.ProjectDir, _ = os.Getwd()
	}

	// Create the project
	crProject := project.New(c.ProjectDir)

	// Generate the plugin.yml file
	pluginYML, err := build.GeneratePluginYML(crProject)
	if err != nil {
		return fmt.Errorf("generating plugin.yml: %w", err)
	}

	// Encode it to stdout
	enc := yaml.NewEncoder(os.Stdout)
	enc.SetIndent(2)
	if err := enc.Encode(pluginYML); err != nil {
		return fmt.Errorf("encoding plugin.yml: %w", err)
	}
	return nil
}
