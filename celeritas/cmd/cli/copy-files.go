package main

import (
	"embed"
	"os"
)

//go:embed templates
var templateFS embed.FS

func copyFileFromTemplate(templatePath, targetFile string) error {
	//TODO check to ensure file does not exist

	data, err := templateFS.ReadFile(templatePath)
	if err != nil {
		exitGracefully(err)
	}

	err = copyDataToFile(data, targetFile)
	if err != nil {
		exitGracefully(err)
	}

	return nil
}

func copyDataToFile(data []byte, to string) error {
	//WriteFile writes data to the named file, creating it if necessary.
	err := os.WriteFile(to, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
