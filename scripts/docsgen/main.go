package main

import (
	"fmt"
	"go/ast"
	"go/doc"
	"os"
	"path/filepath"
	"strings"
	"log"

	"golang.org/x/tools/go/packages"
)

func main() {
	outDir := filepath.Join("docs", "architecture")
	if err := os.MkdirAll(outDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	targets := []string{"internal/core", "internal/features"}

	for _, target := range targets {
		processDirectory(target, outDir)
	}

	log.Println("Documentation generation complete.")
}

func processDirectory(baseDir, outDir string) {
	entries, err := os.ReadDir(baseDir)
	if err != nil {
		log.Printf("Warning: Failed to read %s: %v", baseDir, err)
		return
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		pkgPath := filepath.Join(baseDir, entry.Name())
		generateMarkdownForPackage(pkgPath, baseDir, outDir)
	}
}

func generateMarkdownForPackage(pkgPath, baseDir, outDir string) {
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedFiles |
			packages.NeedSyntax | packages.NeedTypes,
		Dir: pkgPath,
	}

	pkgs, err := packages.Load(cfg, ".")
	if err != nil {
		fmt.Printf(
			"Warning: failed to load package at %s: %v\n",
			pkgPath,
			err,
		)
		return
	}

	if len(pkgs) == 0 {
		return
	}

	for _, pkg := range pkgs {
		if strings.HasSuffix(pkg.Name, "_test") {
			continue
		}

		// Convert *packages.Package back to *ast.Package for doc.New
		//lint:ignore SA1019 ast.Package is deprecated but required by go/doc
		astPkg := &ast.Package{
			Name:  pkg.Name,
			Files: make(map[string]*ast.File),
		}
		for _, f := range pkg.Syntax {
			// pkg.Fset is available in newer versions of x/tools
			fileName := pkg.Fset.Position(f.Pos()).Filename
			astPkg.Files[fileName] = f
		}

		docPkg := doc.New(astPkg, pkgPath, 0)
		generateFile(docPkg, pkgPath, baseDir, outDir)
	}
}

func generateFile(docPkg *doc.Package, pkgPath, baseDir, outDir string) {
	fileName := fmt.Sprintf("%s.md", filepath.Base(pkgPath))
	if strings.Contains(baseDir, "core") {
		fileName = "core_" + fileName
	}

	filePath := filepath.Join(outDir, fileName)
	file, err := os.Create(filePath)
	if err != nil {
		log.Printf("Failed to create file %s: %v", filePath, err)
		return
	}
	defer file.Close()

	layer := "Feature Module"
	if strings.Contains(baseDir, "core") {
		layer = "Core Module"
	}

	fmt.Fprintf(file, "# %s: `%s`\n\n", layer, docPkg.Name)
	if docPkg.Doc != "" {
		fmt.Fprintf(file, "## Overview\n%s\n\n", strings.TrimSpace(docPkg.Doc))
	}

	writeTypes(file, docPkg)
	writeFunctions(file, docPkg.Funcs)

	log.Printf("Generated %s", filePath)
}

func writeTypes(file *os.File, docPkg *doc.Package) {
	if len(docPkg.Types) == 0 {
		return
	}

	fmt.Fprintf(file, "## Types and Interfaces\n\n")
	for _, t := range docPkg.Types {
		fmt.Fprintf(file, "### `%s`\n", t.Name)
		if t.Doc != "" {
			fmt.Fprintf(file, "%s\n\n", strings.TrimSpace(t.Doc))
		}

		// Identify if interface or struct to give context
		kind := "Type"
		if len(t.Methods) > 0 {
			switch t.Decl.Specs[0].(*ast.TypeSpec).Type.(type) {
			case *ast.InterfaceType:
				kind = "Interface"
			case *ast.StructType:
				kind = "Struct"
			}
		} else {
			switch t.Decl.Specs[0].(*ast.TypeSpec).Type.(type) {
			case *ast.InterfaceType:
				kind = "Interface"
			case *ast.StructType:
				kind = "Struct"
			}
		}

		fmt.Fprintf(file, "**Kind**: %s\n\n", kind)

		if len(t.Methods) > 0 {
			fmt.Fprintf(file, "**Methods:**\n")
			for _, m := range t.Methods {
				fmt.Fprintf(file, "- `%s`\n", m.Name)
				if m.Doc != "" {
					fmt.Fprintf(file, "  - *%s*\n", strings.ReplaceAll(strings.TrimSpace(m.Doc), "\n", " "))
				}
			}
			fmt.Fprintf(file, "\n")
		}

		if len(t.Funcs) > 0 {
			fmt.Fprintf(file, "**Constructors/Factory Functions:**\n")
			for _, f := range t.Funcs {
				fmt.Fprintf(file, "- `%s`\n", f.Name)
			}
			fmt.Fprintf(file, "\n")
		}
	}
}

func writeFunctions(file *os.File, funcs []*doc.Func) {
	if len(funcs) == 0 {
		return
	}

	fmt.Fprintf(file, "## Package Level Functions\n\n")
	for _, f := range funcs {
		fmt.Fprintf(file, "### `%s`\n", f.Name)
		if f.Doc != "" {
			fmt.Fprintf(file, "%s\n\n", strings.TrimSpace(f.Doc))
		}
	}
}
