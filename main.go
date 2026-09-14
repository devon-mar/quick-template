package main

import (
	"fmt"
	"log/slog"
	"os"
	"path"
	"text/template"
)

func executeTemplate(funcs template.FuncMap, in string, outPath string) error {
	t, err := template.New(path.Base(in)).Funcs(funcs).ParseFiles(in)
	if err != nil {
		return fmt.Errorf("template.New: %w", err)
	}

	outFile, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("os.Create: %w", err)
	}
	defer outFile.Close()

	if err := t.Execute(outFile, nil); err != nil {
		return fmt.Errorf("execute: %w", err)
	}

	return nil
}

func main() {
	files := os.Args
	if len(files) == 0 {
		slog.Error("no program name")
		os.Exit(1)
	}
	files = os.Args[1:]

	if len(files) == 0 || len(files)%2 != 0 {
		slog.Error("args should be pairs of templates and output files")
		os.Exit(1)
	}

	fileFunc := func(path string) (string, error) {
		b, err := os.ReadFile(path)
		if err != nil {
			return "", err
		}
		return string(b), nil
	}

	templateFuncs := map[string]any{
		"file": fileFunc,
	}

	for i := 0; i < len(files); i += 2 {
		templateFile := files[i]
		outputFile := files[i+1]
		if err := executeTemplate(templateFuncs, templateFile, outputFile); err != nil {
			slog.Error("could not execute template", "template", templateFile, "output", outputFile, "err", err)
			os.Exit(1)
		}
	}
}
