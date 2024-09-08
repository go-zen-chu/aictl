package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

// queryCmd represents the query command
var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "Run query to generative AI",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := validate(args); err != nil {
			return fmt.Errorf("validation in query: %w", err)
		}
		uq := dic.UsecaseQuery()
		query := ""
		if inputStdin {
			stdin, err := io.ReadAll(os.Stdin)
			if err != nil {
				return fmt.Errorf("read stdin: %w", err)
			}
			query = string(stdin)
		} else {
			query = args[0]
		}
		res, err := uq.QueryToOpenAI(query, outputFormat)
		if err != nil {
			return fmt.Errorf("query to openai: %w", err)
		}
		fmt.Printf("%s", res)
		return nil
	},
}

var outputFormat string
var inputStdin bool

func init() {
	rootCmd.AddCommand(queryCmd)

	queryCmd.Flags().StringVarP(&outputFormat, "output", "o", "text", "Output format text or json (default is text)")
	queryCmd.Flags().BoolVarP(&inputStdin, "stdin", "i", false, "Read query from stdin")
}

func validate(args []string) error {
	if !inputStdin && len(args) != 1 {
		return fmt.Errorf("query command requires only 1 argument `query text`")
	}
	if outputFormat == "" {
		return fmt.Errorf("output format is required but got empty")
	}
	if outputFormat != "text" && outputFormat != "json" {
		return fmt.Errorf("output format must be text or json but got: %s", outputFormat)
	}
	return nil
}
