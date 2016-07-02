package args

import (
	"fmt"
	"strings"
)

type Tokenizer struct {
	argv []string
	idx, off int
	arg string
	err error
	noMoreOptions bool
}

func NewTokenizer(argv []string) *Tokenizer {
	return &Tokenizer{
		argv: argv,
	}
}

func (tok *Tokenizer) Next() bool {
	defer func() {
		maybeErr := recover()
		err, isErr := maybeErr.(error)
		if isErr {
			tok.err = err
		}
	}()
	tok.err = nil

	if tok.idx == 0 {
		tok.idx++
	}

	if tok.idx >= len(tok.argv) {
		return false
	}

	tok.arg = tok.argv[tok.idx]

	// Positional argument
	//   arg = positional
	//   off = ^
	if !strings.HasPrefix(tok.arg, "-") || tok.noMoreOptions {
		tok.idx++
		return true
	}

	// No more options marker
	if tok.arg == "--" {
		tok.noMoreOptions = true
		tok.idx++
		return tok.Next()
	}

	// Long option with parameter
	//   arg = --long-option=param
	//   off = ^
	if strings.HasPrefix(tok.arg, "--") && strings.Contains(tok.arg, "=") && tok.off == 0 {
		tok.off = runeIndex(tok.arg, '=')
		tok.arg = tok.arg[:tok.off]
		return true
	}

	// Long option without parameter
	//   arg = --long-option
	//   off = ^
	if strings.HasPrefix(tok.arg, "--") && tok.off == 0 {
		tok.idx++
		return true
	}

	// One or more short options
	//   arg = -xzvf
	//   off = ^^^^^
	//
	//   arg = -Dparameter
	//   off = ^
	//
	//   arg = -o=parameter
	//   off = ^
	if strings.HasPrefix(tok.arg, "-") && !strings.HasPrefix(tok.arg, "--"){
		if tok.off == 0 {
			tok.off++
		}

		// This segment is tricky with state
		c := []rune(tok.arg)[tok.off]
		tok.off++
		if tok.off == len([]rune(tok.arg)) {
			tok.off = 0
			tok.idx++
		}
		tok.arg = fmt.Sprintf("-%c", c)

		return true
	}

	if tok.off != 0 {
		panic(fmt.Errorf("Unexpected parameter '%s'", tok.arg[tok.off:]))
	}

	panic(fmt.Errorf("Failed to process argument '%s'", tok.arg))

	return false
}

func (tok *Tokenizer) Arg() string {
	return tok.arg
}

func (tok *Tokenizer) IsOption() bool {
	return strings.HasPrefix(tok.arg, "-") && !tok.noMoreOptions
}

func (tok *Tokenizer) Err() error {
	return tok.err
}

func (tok *Tokenizer) TakeParameter() (string, error) {
	if tok.idx == len(tok.argv) {
		return "", fmt.Errorf("Missing parameter to %s option", tok.arg)
	}

	param := tok.argv[tok.idx]

	if tok.off == 0 {
		tok.idx++
		return param, nil
	}

	if runeIndex(param, '=') == tok.off {
		param = string([]rune(param)[tok.off + 1:])
		tok.idx++
		tok.off = 0
		return param, nil
	}

	if strings.HasPrefix(param, "-") && !strings.HasPrefix(param, "--") && tok.off >= 2 {
		param = string([]rune(param)[tok.off:])
		tok.idx++
		tok.off = 0
		return param, nil
	}

	return "", fmt.Errorf("Failed to parse parameter '%s'", param[tok.off:])
}

func (tok *Tokenizer) Remainder() ([]string, bool) {
	return tok.argv[tok.idx:], tok.noMoreOptions
}

// strings.Index works by bytes, not runes
func runeIndex(s string, r rune) int {
	for i, v := range []rune(s) {
		if r == v { return i }
	}
	return -1
}

