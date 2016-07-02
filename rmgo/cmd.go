package main

import (
	"go/token"
	"go/ast"
	"fmt"
	"os"
	"io"
)

type cmdContext struct{
	FileSet *token.FileSet
	Ast ast.Node 
	DeclMap map[string]ast.Decl
	Args []string
	Stdin io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

type cmdStatus int

type cmdListItem struct {
	name string
	fn func(*cmdContext)
}

var cmdList = map[string]*cmdListItem{}

func cmdRegister(name string, fn func(*cmdContext)) {
	cmdList[name] = &cmdListItem{ name: name, fn: fn }
}

func cmdNewContext(astNode ast.Node, fset *token.FileSet, declMap map[string]ast.Decl) (*cmdContext) {
	return &cmdContext{
		FileSet: fset,
		Ast: astNode,
		DeclMap: declMap,
		Stdin: os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
}

func (cmd *cmdContext) Run(argv ...string) (status int) {

	cmd.Args = argv

	command, commandExists := cmdList[argv[0]]
	if !commandExists {
		fmt.Fprintf(cmd.Stderr, "Command does not exist.\n")
		return 1
	}

	defer func() {
		r := recover()
		if r == nil {
			return
		}
		switch r.(type) {
		case cmdStatus:
			status = r.(int)
		default:
			panic(r)
		}
	}()
	command.fn(cmd)
		
	return status
}

func (cmd *cmdContext) Print(args... interface{}) (int, error) {
	return fmt.Fprint(cmd.Stdout, args...)
}

func (cmd *cmdContext) Println(args... interface{}) (int, error) {
	return fmt.Fprintln(cmd.Stdout, args...)
}

func (cmd *cmdContext) Printf(format string, args... interface{}) (int, error) {
	return fmt.Fprintf(cmd.Stdout, format, args...)
}

func (cmd *cmdContext) Fatal(args... interface{}) {
	fmt.Fprint(cmd.Stderr, args...)
	cmd.Exit(1)
}

func (cmd *cmdContext) Fatalln(args... interface{}) {
	fmt.Fprintln(cmd.Stderr, args...)
	cmd.Exit(1)
}

func (cmd *cmdContext) Fatalf(format string, args... interface{}) {
	fmt.Fprintf(cmd.Stderr, format, args...)
	cmd.Exit(1)
}

func (cmd *cmdContext) Exit(code int) {
	panic(cmdStatus(code))
}

