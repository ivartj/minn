package main

import (
	"github.com/ivartj/minn/args"
	"go/ast"
	"container/list"
	"io"
	"os"
	"fmt"
	"bytes"
	"bufio"
)

func init() {
	cmdRegister("remove", removeCmd)
}

func removeCmdUsage(w io.Writer) {
	fmt.Fprintf(w, "Usage: %s remove <identifier> ...\n", mainProgramName)
}

func removeCmdArgs(cmd *cmdContext) (idents []string) {

	tok := args.NewTokenizer(cmd.Args)
	idents = []string{}

	for tok.Next() {

		if tok.IsOption() {

			switch tok.Arg() {

			case "-h", "--help":
				removeCmdUsage(cmd.Stdout)
				cmd.Exit(0)

			default:
				cmd.Fatalf("Unrecognized option, '%s'.\n", tok.Arg())

			}

		} else {
			idents = append(idents, tok.Arg())
		}

	}

	if tok.Err() != nil {
		cmd.Fatalf("Error occurred on processing command-line options: %s.\n", tok.Err().Error())
	}

	if len(idents) == 0 {
		removeCmdUsage(cmd.Stderr)
		cmd.Exit(1)
	}

	return idents
}

func removeCmdGetDeclsByFilename(cmd *cmdContext, idents []string) map[string][]ast.Decl {

	m := map[string][]ast.Decl{}

	for _, ident := range idents {

		decl, declExists := cmd.DeclMap[ident]
		if !declExists {
			cmd.Fatalf("Did not find declaration for '%s'.\n", ident)
		}

		filename := cmd.FileSet.Position(decl.Pos()).Filename
		declList := m[filename]
		if declList == nil {
			declList = []ast.Decl{}
		}
		m[filename] = append(declList, decl)
	}

	return m
}

type removeCmdRange struct{
	pos, end int
}

func removeCmdGetLineRanges(cmd *cmdContext, decls []ast.Decl) []removeCmdRange {

	set := make([]removeCmdRange, len(decls))

	for i, decl := range decls {
		pos := cmd.FileSet.Position(decl.Pos()).Line
		end := cmd.FileSet.Position(decl.End()).Line
		set[i] = removeCmdRange{pos, end}
	}

	ls := list.New()

	for i, r := range set {

		if i == 0 {
			ls.PushFront(r)
			continue
		}

		inserted := false
		for el := ls.Front(); el != nil && inserted == false; {

			v := el.Value.(removeCmdRange)

			switch {

			// v is before r
			case r.pos > v.end:
				el = el.Next()

			// r is before v
			case r.end < v.pos:
				ls.InsertBefore(r, el)
				inserted = true

			// Overlap
			default:
				if v.pos < r.pos { r.pos = v.pos }
				if v.end > r.end { r.end = v.end }
				next := el.Next()
				ls.Remove(el)
				el = next
			}
		}

		if inserted == false {
			ls.PushBack(r)
		}
	}

	set = set[0:ls.Len()]
	i := 0
	for el := ls.Front(); el != nil; el = el.Next() {
		set[i] = el.Value.(removeCmdRange)
		i++
	}

	return set
}

func removeCmdRemoveLineRanges(cmd *cmdContext, filename string, lineRanges []removeCmdRange) {

	file, err := os.OpenFile(filename, os.O_RDWR, 0)
	if err != nil {
		cmd.Fatalf("Error occurred on opening file '%s': %s.\n", filename, err.Error())
	}
	defer file.Close()

	buf := bytes.NewBuffer([]byte{})
	br := bufio.NewReader(file)

	for nline := 1;; nline++ {

		line, err := br.ReadBytes('\n')
		if err != nil && err != io.EOF {
			cmd.Fatalf("Error on reading line from %s: %s.\n", filename, err.Error())
		}

		if len(line) == 0 {
			break
		}

		for len(lineRanges) > 0 && nline > lineRanges[0].end {
			lineRanges = lineRanges[1:]
		}

		if len(lineRanges) != 0 {

			r := lineRanges[0]

			switch {

			// line is before current line range, no effect
			case nline < r.pos:

			// line is within current line range, skipped
			case nline >= r.pos && nline <= r.end:
				continue

			}
			
		}

		_, err = buf.Write(line)
		if err != nil {
			cmd.Fatalf("Error on writing to buffer: %s.\n", err.Error())
		}
	}

	_, err = file.Seek(0, 0)
	if err != nil {
		cmd.Fatalf("Failed to seek to beginning of %s: %s.\n", filename, err.Error())
	}

	err = file.Truncate(int64(buf.Len()))
	if err != nil {
		cmd.Fatalf("Failed to truncate %s: %s.\n", filename, err.Error())
	}

	_, err = buf.WriteTo(file)
	if err != nil {
		cmd.Fatalf("Error occurred on writing buffer to %s: %s.\n", filename, err.Error())
	}
}

func removeCmd(cmd *cmdContext) {

	idents := removeCmdArgs(cmd)

	declsByFilename := removeCmdGetDeclsByFilename(cmd, idents)

	for filename, decls := range declsByFilename {
		lineRanges := removeCmdGetLineRanges(cmd, decls)
		removeCmdRemoveLineRanges(cmd, filename, lineRanges)
	}
}

