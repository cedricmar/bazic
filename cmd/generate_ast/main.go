package main

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
	"strings"
)

func main() {

	if len(os.Args) != 1 {
		fmt.Println("Usage: generate_ast")
		os.Exit(64)
	}

	dir := "./pkg/ast"

	defineAst(dir, "Expr", []string{
		"Binary   : left Expr, operator scanner.Token, right Expr",
		"Grouping : expression Expr",
		"Literal  : value interface{}",
		"Unary    : operator scanner.Token, right Expr",
	})
}

func defineAst(outputDir, baseName string, types []string) error {

	path := outputDir + "/" + strings.ToLower(baseName) + ".go"

	buf := bytes.Buffer{}

	buf.WriteString("package ast\n")
	buf.WriteString("\n")
	buf.WriteString("import \"github.com/cedricmar/bazic/pkg/scanner\"\n")
	buf.WriteString("\n")
	buf.WriteString("// " + baseName + " is a type for the AST\n")
	buf.WriteString("type " + baseName + " string\n")

	for _, t := range types {
		c := strings.Trim(strings.Split(t, ":")[0], " ")
		fs := strings.Trim(strings.Split(t, ":")[1], " ")

		defineType(&buf, c, fs)
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	// Format
	src, err := format.Source(buf.Bytes())
	if err != nil {
		return err
	}

	_, err = f.Write(src)
	if err != nil {
		return err
	}

	return nil
}

func defineType(buf *bytes.Buffer, class string, fields string) {
	buf.WriteString("\n")

	// Declare the struct type
	buf.WriteString("// " + class + " is a node of the AST\n")
	buf.WriteString("type " + class + " struct {\n")

	// Write fields for the struct type
	fs := strings.Split(fields, ", ")
	for _, f := range fs {
		buf.WriteString("    " + f + "\n")
	}

	buf.WriteString("}\n")
	buf.WriteString("\n")

	// Declare the constructor
	buf.WriteString("// New" + class + " returns a new node of type " + class + "\n")
	buf.WriteString("func New" + class + "(" + fields + ") " + class + " {\n")

	// Return the type
	buf.WriteString("    return " + class + "{\n")

	for _, f := range fs {
		v := strings.Split(f, " ")[0]
		buf.WriteString("        " + v + ": " + v + ",\n")
	}

	buf.WriteString("    }\n")

	buf.WriteString("}\n")
}
