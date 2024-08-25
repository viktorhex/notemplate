package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	// expect 2 or less args
	if len(os.Args) > 3 {
		fmt.Println("too many args")
		os.Exit(1)
	}

	// arg 1
	var templateName string = ""
	if len(os.Args) > 1 && os.Args[1] != "_" {
		templateName = os.Args[1]
	}

	// arg 2
	var fileSuffix string = ""
	if len(os.Args) == 3 {
		fileSuffix = os.Args[2]
	}

	// output dir
	var entriesDir string = "notemplates"
	if templateName != "" {
		entriesDir = templateName
	}

	// try to load content
	var content string
	var loadTemplateErr error
	if templateName != "" {
		content, loadTemplateErr = loadTemplate(templateName)
		if loadTemplateErr != nil {
			fmt.Printf("Error loading template: %v\n", loadTemplateErr)
			os.Exit(1)
		}
	} else {
		content = "# Notemplate entry\n"
	}
	err := os.MkdirAll(entriesDir, 0755)
	if err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		return
	}

	// create file name
	currentDate := time.Now().Format("2006-01-02")
	n := 0
	var foldername string
	println(fileSuffix)
	if fileSuffix != "" && fileSuffix != "_" {
		fileSuffix = "-" + fileSuffix
	}

	// increment n until a new file name is available
	for {
		foldername = fmt.Sprintf("%s-entry-%d%s", currentDate, n, fileSuffix)
		foldernameNosuffix := fmt.Sprintf("%s-entry-%d", currentDate, n)
		fullPath := filepath.Join(entriesDir, foldername)
		fullPathNosuffix := filepath.Join(entriesDir, foldernameNosuffix)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			if _, err := os.Stat(fullPathNosuffix); os.IsNotExist(err) {
				break // break only if neither file nor file with suffix exist
			}
		}
		n++
	}

	if err := os.Mkdir(filepath.Join(entriesDir, foldername), 0755); err != nil {
		println("%s", err.Error())
		os.Exit(1)
	}

	if templateName == "job_applications" {
		filenames := []string{
			"info.toml",
			"events.toml",
		}

		for _, filename := range filenames {
			dirPath := filepath.Join(entriesDir, foldername)
			tmplName := filename
			var text string
			if filename == "info.toml" {
				text = content // main template for this already loaded
			} else {
				text, loadTemplateErr = loadTemplate(tmplName)
				if loadTemplateErr != nil {
					fmt.Printf("Error loading template: %v\n", loadTemplateErr)
					os.Exit(1)
				}
			}
			createFile(dirPath, filename, text)
		}
	} else {
		dirPath := filepath.Join(entriesDir, foldername)
		createFile(dirPath, "info.toml", content)
	}
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

func createFile(entriesDir string, filename string, content string) {
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
