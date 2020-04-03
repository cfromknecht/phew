package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"os"

	"github.com/cfromknecht/phew"
	"github.com/spf13/cobra"
)

var (
	ErrNumEncodeArgs = errors.New("encode only acceepts zero or one arguments")
)

// encodeCmd represents the encode command
var encodeCmd = &cobra.Command{
	Use:   "encode [data]",
	Short: "Encode data passed as either a single argument or from stdin.",
	Long: `Encode accepts zero or one arguments. 
	
When no arguments are passed, encode reads data from stdin and writes the
encoded output to stdout.

When one argument is passed, encode will write the encoded argument to stdout.

Passing in more than one argument will fail.`,
	Run: func(cmd *cobra.Command, args []string) {
		data := getInputData(args, ErrNumEncodeArgs)

		var b bytes.Buffer
		w := phew.NewWriter(&b)
		write(w, data)
		err := w.Close()
		if err != nil {
			fatal(fmt.Errorf("unable to close phew writer: %v", err))
		}

		write(os.Stdout, b.Bytes())
		write(os.Stdout, newlineBytes)
	},
}

func init() {
	rootCmd.AddCommand(encodeCmd)
}
