package main

import (
	"github.com/ivartj/minn/args"
	"fmt"
	"io"
	"os"
	"go/ast"
	"go/token"
	"go/parser"
	"go/printer"
	"bytes"
)

const (
	mainProgramName		= "rmgo"
	mainProgramVersion	= "0.1-SNAPSHOT"
)

var (
	mainConfSource		= "."
)

func mainUsage(w io.Writer) {
	fmt.Fprintf(w, "Usage: %s [ -s <source> ] <command> ...\n", mainProgramName)
}

func mainArgs() (argv []string) {

	tok := args.NewTokenizer(os.Args)
	commandGiven := false
	command := ""

	for tok.Next() {

		if tok.IsOption() {

			switch tok.Arg() {

			case "-h", "--help":
				mainUsage(os.Stdout)
				os.Exit(0)

			case "--version":
				fmt.Printf("%s version %s\n", mainProgramName, mainProgramVersion)
				os.Exit(0)

			case "-s", "--source":
				param, err := tok.TakeParameter()
				if err != nil {
					fmt.Fprintf(os.Stderr, "Failed to get parameter to %s option: %s.\n", err.Error())
					os.Exit(1)
				}
				mainConfSource = param

			default:
				fmt.Fprintf(os.Stderr, "Unrecognized option, '%s'.\n", tok.Arg())
				os.Exit(1)

			}

		} else {
			commandGiven = true
			command = tok.Arg()
			break
		}

	}

	if tok.Err() != nil {
		fmt.Fprintf(os.Stderr, "Error occurred on processing command-line arguments: %s.\n", tok.Err().Error())
		os.Exit(1)
	}

	if !commandGiven {
		mainUsage(os.Stderr)
		os.Exit(1)
	}

	argv, noMoreOptions := tok.Remainder()
	if noMoreOptions {
		fmt.Fprintf(os.Stderr, "No-more-options marker (--) not allowed before subcommand.\n")
		os.Exit(1)
	}

	return append([]string{command}, argv...)
}

func mainParse() (ast.Node, *token.FileSet) {

	fileInfo, err := os.Stat(mainConfSource)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to stat %s: %s.\n", mainConfSource, err.Error())
		os.Exit(1)
	}

	fset := token.NewFileSet()

	switch {
	case fileInfo.IsDir():
		pkgmap, err := parser.ParseDir(fset, mainConfSource, nil, parser.AllErrors)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error on parsing files in %s: %s.\n", mainConfSource, err.Error())
			os.Exit(1)
		}

		if len(pkgmap) != 1 {
			fmt.Fprintf(os.Stderr, "Not 1 but %d packages in %s.\n", len(pkgmap), mainConfSource)
			os.Exit(1)
		}

		var pkg *ast.Package
		for _, v := range pkgmap { pkg = v; break }

		return pkg, fset

	case fileInfo.Mode().IsRegular():
		file, err := parser.ParseFile(fset, mainConfSource, nil, parser.AllErrors)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error on parsing %s: %s.\n", mainConfSource, err.Error())
			os.Exit(1)
		}
		return file, fset

	default:
		fmt.Fprintf(os.Stderr, "%s does not appear to be a directory nor a regular file.\n", mainConfSource)
		os.Exit(1)
	}

	return nil, nil
}

func mainMapSymbols(fset *token.FileSet, declMap map[string]ast.Decl, node ast.Node) {

	switch n := node.(type) {

	case *ast.Package:
		for _, file := range n.Files {
			mainMapSymbols(fset, declMap, file)
		}

	case *ast.File:
		for _, decl := range n.Decls {
			mainMapSymbols(fset, declMap, decl)
		}

	case *ast.FuncDecl:
		buf := bytes.NewBuffer([]byte{})
		if n.Recv != nil && n.Recv.List != nil {
			fmt.Fprint(buf, "(")
			for i, field := range n.Recv.List {
				if i != 0 {
					fmt.Fprint(buf, ",")
				}
				err := printer.Fprint(buf, fset, field.Type)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error occurred on printing receiver type: %s.\n", err.Error())
					os.Exit(1)
				}
			}
			fmt.Fprint(buf, ").")
		} else {
			if n.Name.Name == "init" {
				return
			}
		}
		fmt.Fprint(buf, n.Name.Name)
		declMap[buf.String()] = n

	case *ast.GenDecl:
		for _, spec := range n.Specs {
			switch s := spec.(type) {
			case *ast.ValueSpec:
				for _, id := range s.Names {
					// TODO: Find more specific declaration node
					declMap[id.Name] = n
				}
			case *ast.TypeSpec:
				declMap[s.Name.Name] = n
			}
		}

	default:
		fmt.Fprintf(os.Stderr, "Unexpected ast.Node type %T.\n", n)
		os.Exit(1)
	}
}

func main() {
	argv := mainArgs()
	astNode, fset := mainParse()
	declMap := map[string]ast.Decl{}
	mainMapSymbols(fset, declMap, astNode)
	cmd := cmdNewContext(astNode, fset, declMap)
	os.Exit(cmd.Run(argv...))
}

