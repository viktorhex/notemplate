package main

import (
	"embed"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	fmt.Println("notemplate version 1.0")
	var template string
	var suffix string
	const (
		templateDefault = ""
		templateUsage   = "the notes template to create"
		suffixDefault   = ""
		suffixUsage     = "text to add to end of created file (or folder)"
	)
	flag.StringVar(&template, "template", templateDefault, templateUsage)
	flag.StringVar(&template, "t", templateDefault, templateUsage+" (shorthand)")
	flag.StringVar(&suffix, "suffix", suffixDefault, suffixUsage)
	flag.StringVar(&suffix, "s", suffixDefault, suffixUsage+" (shorthand)")

	flag.Parse()

	if template == "" {
		fmt.Println("ERR: --template is required (shorthand: -t)")
		os.Exit(1)
	}

	createEntryParams := CreateEntryParams{template, suffix}
	create_entry(createEntryParams)
}

type CreateEntryParams struct {
	template, suffix string
}

//go:embed templates/*.toml
var templates embed.FS

func create_entry(p CreateEntryParams) {

	// output dir
	entriesDir := "notemplates"
	if p.template != "" {
		entriesDir = p.template
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	templatesRoot := homeDir + "/notemplate/templates"
	templateFiles := []string{
		"notes.toml",
		"star_entries.toml",
		"research_notes.toml",
	}

	if err := CopyEmbeddedFilesToFolder(templatesRoot, templateFiles); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println("TOML files copied successfully")

	// try to load content
	content := ""
	var loadTemplateErr error
	if p.template != "" {
		content, loadTemplateErr = loadTemplate(p.template)
		if loadTemplateErr != nil {
			fmt.Printf("Error loading template: %v\n", loadTemplateErr)
			os.Exit(1)
		}
	} else {
		content = "# Notemplate entry\n"
	}

	docsRoot := homeDir + "/notemplate/documents"

	err1 := os.MkdirAll(docsRoot, 0755)
	if err1 != nil {
		fmt.Printf("Error creating directory: %v\n", err1)
		return
	}

	err2 := os.MkdirAll(path.Join(docsRoot, entriesDir), 0755)
	if err2 != nil {
		fmt.Printf("Error creating directory: %v\n", err1)
		return
	}

	// create file name
	currentDate := time.Now().Format("2006-01-02")
	n := 0
	var foldername string
	println(p.suffix)
	if p.suffix != "" && p.suffix != "_" {
		p.suffix = "-" + p.suffix
	}

	// increment n until a new file name is available
	for {
		foldername = fmt.Sprintf("%s-entry-%d%s", currentDate, n, p.suffix)
		foldernameNosuffix := fmt.Sprintf("%s-entry-%d", currentDate, n)
		fullPath := filepath.Join(docsRoot, entriesDir, foldername)
		fullPathNosuffix := filepath.Join(docsRoot, entriesDir, foldernameNosuffix)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			if _, err := os.Stat(fullPathNosuffix); os.IsNotExist(err) {
				break // break only if neither file nor file with suffix exist
			}
		}
		n++
	}

	if err := os.Mkdir(filepath.Join(docsRoot, entriesDir, foldername), 0755); err != nil {
		println("%s", err.Error())
		os.Exit(1)
	}

	if p.template == "job_applications" {
		filenames := []string{
			"info.toml",
			"events.toml",
		}

		for _, filename := range filenames {
			dirPath := filepath.Join(docsRoot, entriesDir, foldername)
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
		dirPath := filepath.Join(docsRoot, entriesDir, foldername)
		createFile(dirPath, "info.toml", content)
	}
}

func loadTemplate(templateName string) (string, error) {

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	templatesRoot := homeDir + "/notemplate/templates"

	templatePath := filepath.Join(templatesRoot, templateName)
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

func CopyEmbeddedFilesToFolder(destFolder string, templateNames []string) error {
	if err := os.MkdirAll(destFolder, 0755); err != nil {
		return fmt.Errorf("failed to create destination folder: %v", err)
	}

	for _, templateName := range templateNames {
		destPath := filepath.Join(destFolder, templateName)

		if _, err := os.Stat(destPath); !os.IsNotExist(err) {
			continue
		}

		data, err := templates.ReadFile(filepath.Join("templates", templateName))
		if err != nil {
			return fmt.Errorf("failed to read embedded TOML file %s: %v", templateName, err)
		}

		destFile, err := os.Create(destPath)
		if err != nil {
			return fmt.Errorf("failed to create destination file %s: %v", destPath, err)
		}
		defer destFile.Close()

		if _, err := destFile.Write(data); err != nil {
			return fmt.Errorf("failed to write to destination file %s: %v", destPath, err)
		}
	}

	return nil
}
