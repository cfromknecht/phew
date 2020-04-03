package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/cfromknecht/phew"
	"github.com/spf13/cobra"
)

var (
	ErrNumDecodeArgs = errors.New("decode only accepts zero or one arguments")
)

// decodeCmd represents the decode command
var decodeCmd = &cobra.Command{
	Use:   "decode [data]",
	Short: "Decode data passed as either a single argument or from stdin.",
	Long: `Decode accepts zero or one arguments.

When no arguments are passed, decode reads data from stdin and writes the
encoded output to stdout.

When one argument is passed, encode will write the decoded argument to stdout.

Passing in more than one argument will fail.`,
	Run: func(cmd *cobra.Command, args []string) {
		data := getInputData(args, ErrNumDecodeArgs)

		r, err := phew.NewReader(bytes.NewReader(data))
		if err != nil {
			fatal(fmt.Errorf("unable to create phew reader: %v", err))
		}

		plaintext, err := ioutil.ReadAll(r)
		if err != nil {
			fatal(fmt.Errorf("unable to decode phewed data: %v", err))
		}

		write(os.Stdout, plaintext)
		write(os.Stdout, newlineBytes)
	},
}

func init() {
	rootCmd.AddCommand(decodeCmd)
}
