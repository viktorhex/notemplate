package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	if len(os.Args) > 3 {
		fmt.Println("too many args")
		os.Exit(1)
	}
	var templateName string = ""
	if len(os.Args) == 2 {
		templateName = os.Args[1]
	}
	var fileSuffix string = ""
	if len(os.Args) == 3 {
		fileSuffix = os.Args[2]
	}
	var entriesDir string = "notemplates"
	if templateName != "" {
		entriesDir = templateName
	}
	var content string
	var loadTemplateErr error
	if templateName != "" {
		content, loadTemplateErr = loadTemplate(templateName)
		if loadTemplateErr != nil {
			fmt.Printf("Error loading template: %v\n", loadTemplateErr)
			return
		}
	} else {
		content = "# Notemplate entry\n"
	}
	err := os.MkdirAll(entriesDir, 0755)
	if err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		return
	}

	currentDate := time.Now().Format("2006-01-02")
	n := 0
	var filename string
	for {
		if fileSuffix != "" {
			fileSuffix = "-" + fileSuffix
		}
		filename = fmt.Sprintf("%s-entry-%d%s.toml", currentDate, n, fileSuffix)
		fullPath := filepath.Join(entriesDir, filename)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			break
		}
		n++
	}

	fullPath := filepath.Join(entriesDir, filename)
	file, err := os.Create(fullPath)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		fmt.Printf("Error writing to file: %v\n", err)
		return
	}
	fmt.Printf("Created file: %s\n", fullPath)
}

func loadTemplate(templateName string) (string, error) {
	templatePath := filepath.Join("templates", templateName)
	if !strings.HasSuffix(templatePath, ".toml") {
		templatePath += ".toml"
	}
	content, err := os.ReadFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to read template file: %v", err)
	}
	return string(content), nil
}
