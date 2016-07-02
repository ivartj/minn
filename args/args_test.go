package args

import (
	"testing"
)

func TestTokenizer(t *testing.T) {
	argv := []string{
		"cmd",
		"-h",
		"--version",
		"-xzvf",
		"pos1",
		"--output=filename",
		"-i=input",
		"-c", "config.cfg",
		"--",
		"--pos2",
	}

	tok := NewTokenizer(argv)

	var (
		optHelp, optVersion bool
		optX, optZ, optV, optF bool
		optOutput, optInput bool
		optConfig bool
		pos1, pos2 bool
	)

	for tok.Next() {
		arg := tok.Arg()
		isOption := tok.IsOption()

		switch {
		case !isOption:
			switch arg {
			case "pos1":
				pos1 = true
			case "--pos2":
				pos2 = true
			default:
				t.Fatalf("Unexpected positional '%s'.\n", arg)
			}

		case isOption:
			switch arg {
			case "-h", "--help":
				optHelp = true

			case "--version":
				optVersion = true

			case "-x": optX = true
			case "-z": optZ = true
			case "-v": optV = true
			case "-f": optF = true

			case "-o", "--output":
				optOutput = true
				param, err := tok.TakeParameter()
				if err != nil {
					t.Fatalf("Error on calling Tokenizer.TakeParameter: %s\n", err.Error())
				}
				if param != "filename" {
					t.Fatalf("'%s' does not match the expected parameter 'filename'.\n", param)
				}

			case "-i", "--input":
				optInput = true
				param, err := tok.TakeParameter()
				if err != nil {
					t.Fatalf("Error on calling Tokenizer.TakeParameter: %s\n", err.Error())
				}
				if param != "input" {
					t.Fatalf("'%s' does not match the expected parameter 'input'.\n", param)
				}

			case "-c", "--config":
				optConfig = true
				param, err := tok.TakeParameter()
				if err != nil {
					t.Fatalf("Error on calling Tokenizer.TakeParameter: %s\n", err.Error())
				}
				if param != "config.cfg" {
					t.Fatalf("'%s' does not match the expected parameter 'config.cfg'.\n", param)
				}

			default:
				t.Fatalf("Unexpected option '%s'.\n", arg)
			}
		}
	}

	err := tok.Err()
	if err != nil {
		t.Fatalf("Unexpected error after Tokenizer.Next(): %s\n", err.Error())
	}

	if !optHelp { t.Error("Expected help option not encountered.\n") }
	if !optVersion { t.Error("Expected version option not encountered.\n") }
	if !optX { t.Error("Expected x option not encountered.\n") }
	if !optZ { t.Error("Expected z option not encountered.\n") }
	if !optV { t.Error("Expected v option not encountered.\n") }
	if !optF { t.Error("Expected f option not encountered.\n") }
	if !optOutput { t.Error("Expected output option not encountered.\n") }
	if !optInput { t.Error("Expected input option not encountered.\n") }
	if !optConfig { t.Error("Expected config option not encountered.\n") }
	if !pos1 { t.Error("Expected first positional not encountered.\n") }
	if !pos2 { t.Error("Expected second positional not encountered.\n") }
}

