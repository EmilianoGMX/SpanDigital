package main

import (
	"io"
	"os"
	"slices"
	"testing"
)

func TestStartTournament(t *testing.T) {
	var err error
	var stdOutPipe *os.File
	var r *os.File
	var w *os.File
	var out []byte
	var expected []byte

	stdOutPipe = os.Stdout //Assign all the os stdOut to the variable

	r, w, err = os.Pipe()
	if err != nil {
		t.Error("Testing error: ", err)
	}

	os.Stdout = w //Assign the stdout pointer to the w writter variable

	//We will use the default file because we know what is the output ("matches_expected.txt")
	StartTournament("matches.txt") //All the output printed will be written to w

	w.Close() //Once we got the output we close the writter

	out, err = io.ReadAll(r) //Now we use the reader to get the output
	if err != nil {
		t.Error("Testing error: ", err)
	}

	os.Stdout = stdOutPipe //Reassing the stdout to the default's

	//Usually there is a newline char at the end of the output, so we remove it
	if out[len(out)-1] == '\n' {
		out = out[:len(out)-1]
	}

	//The file is read
	expected, err = os.ReadFile("matches_expected.txt")
	if err != nil {
		t.Error("Testing error: ", err)
	}

	//Compare the slices
	if !slices.Equal(expected, out) {
		t.Errorf("Wrong expected output.\nExpected:\n%v\n\nGot:\n%v\n\n", string(expected), string(out))
	}
}
