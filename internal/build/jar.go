package build

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/customrealms/cli/internal/pluginyml"
	"github.com/customrealms/cli/internal/project"
	"gopkg.in/yaml.v3"
)

type JarAction struct {
	Project     project.Project
	JarTemplate JarTemplate
	BundleFile  string
	OutputFile  string
}

func (a *JarAction) Run(ctx context.Context) error {
	fmt.Println("Downloading runtime")

	// Get the reader of the Jar file
	jarReader, err := a.JarTemplate.Jar()
	if err != nil {
		return err
	}
	defer jarReader.Close()

	// Copy the jar file to a buffer
	var jarTemplateBuf bytes.Buffer
	if _, err := io.Copy(&jarTemplateBuf, jarReader); err != nil {
		return err
	}

	// Make sure the directory above the output file exists
	if err := os.MkdirAll(filepath.Dir(a.OutputFile), 0777); err != nil {
		return err
	}

	// Open the output file for the final JAR
	file, err := os.Create(a.OutputFile)
	if err != nil {
		fmt.Println("Couldn't find output file during the bundling process, check the error below. Directory: '" + a.OutputFile + "'")
		return err
	}
	defer file.Close()

	// Create a reader for the plugin source code
	pluginCode, err := os.Open(a.BundleFile)
	if err != nil {
		return err
	}
	defer pluginCode.Close()

	// Generate the plugin.yml file for the project
	pluginYML, err := GeneratePluginYML(a.Project)
	if err != nil {
		return fmt.Errorf("generating plugin.yml: %w", err)
	}

	// Produce the final JAR file
	if err := WriteJarFile(
		file,
		jarTemplateBuf.Bytes(),
		pluginCode,
		pluginYML,
	); err != nil {
		return err
	}

	fmt.Println("Plugin Destination:", a.OutputFile)

	return nil

}

func WriteJarFile(
	writer io.Writer,
	templateJarData []byte,
	pluginSourceCode io.Reader,
	pluginYML *pluginyml.Plugin,
) error {
	fmt.Println("Packaging...")

	// Create the ZIP writer
	zw := zip.NewWriter(writer)
	defer zw.Close()

	// Create the ZIP reader from the base YML
	templateJarReader := bytes.NewReader(templateJarData)
	zr, err := zip.NewReader(templateJarReader, int64(len(templateJarData)))
	if err != nil {
		return err
	}

	// Copy all the files back to the jar file
	for _, f := range zr.File {

		// Skip some files
		if f.Name == "plugin.js" || f.Name == "plugin.yml" {
			continue
		}

		// Copy the rest
		if err := zw.Copy(f); err != nil {
			return err
		}

	}

	// Write the plugin code to the jar
	codeFile, err := zw.Create("plugin.js")
	if err != nil {
		return err
	}
	if _, err := io.Copy(codeFile, pluginSourceCode); err != nil {
		return err
	}

	// Write the plugin YML file to the jar
	ymlFile, err := zw.Create("plugin.yml")
	if err != nil {
		return err
	}
	enc := yaml.NewEncoder(ymlFile)
	enc.SetIndent(2)
	if err := enc.Encode(pluginYML); err != nil {
		return fmt.Errorf("encoding plugin.yml: %w", err)
	}

	// We're done, no errors
	return nil

}
