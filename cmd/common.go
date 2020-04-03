package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
)

var (
	newlineBytes = []byte("\n")
)

func getInputData(args []string, argsErr error) []byte {
	switch len(args) {
	case 0:
		return readStdin()

	case 1:
		return []byte(args[0])

	default:
		fatal(argsErr)
		return nil
	}
}

func write(w io.Writer, b []byte) {
	_, err := w.Write(b)
	if err != nil {
		fatal(err)
	}
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

func readStdin() []byte {
	reader := bufio.NewReader(os.Stdin)
	data, err := reader.ReadBytes('\n')
	if err != nil {
		fatal(fmt.Errorf("unable to read from stdin: %v", err))
	}

	return bytes.TrimSuffix(data, newlineBytes)
}
