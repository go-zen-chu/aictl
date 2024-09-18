package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

func NewQueryCmd(cmdReq CommandRequirements) *cobra.Command {
	const defaultOutputFormat = "text"
	const defaultResponseLanguage = "English"
	const defaultInputStdin = false
	var defaultTextFilePaths = []string{}

	var outputFormat string
	var responseLanguage string
	var inputStdin bool
	var textFilePaths []string

	queryCmd := &cobra.Command{
		Use:   "query",
		Short: "Run query to generative AI",
		Long: `A longer description that spans multiple lines and likely contains examples
	and usage of using your command. For example:
	
	Cobra is a CLI library for Go that empowers applications.
	This application is a tool to generate the needed files
	to quickly create a Cobra application.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			validate := func(args []string) error {
				if !inputStdin && len(args) != 1 {
					return fmt.Errorf("query command requires only 1 argument `query text`")
				}
				if inputStdin && len(args) != 0 {
					return fmt.Errorf("query command with stdin cannot have a `query text`")
				}
				if outputFormat == "" {
					return fmt.Errorf("output format is required but got empty")
				}
				of := outputFormat
				if of != "text" && of != "json" {
					return fmt.Errorf("output format must be text or json but got: %s", of)
				}
				return nil
			}
			if err := validate(args); err != nil {
				return fmt.Errorf("validation in query: %w", err)
			}
			var query string
			if inputStdin {
				stdin, err := io.ReadAll(os.Stdin)
				if err != nil {
					return fmt.Errorf("read stdin: %w", err)
				}
				query = string(stdin)
			} else {
				query = args[0]
			}
			uq := cmdReq.UsecaseQuery()
			res, err := uq.QueryToOpenAI(
				query,
				outputFormat,
				responseLanguage,
				textFilePaths,
			)
			if err != nil {
				return fmt.Errorf("query to openai: %w", err)
			}
			// Print the response to **stdout**. All the logs are printed to **stderr**
			fmt.Printf("%s", res)
			return nil
		},
	}
	queryCmd.Flags().StringVarP(&outputFormat,
		"output", "o",
		defaultOutputFormat,
		"Output format text or json (default is text)",
	)
	queryCmd.Flags().StringVarP(&responseLanguage,
		"language", "l",
		defaultResponseLanguage,
		"Which language you want to get response (default is English)",
	)
	queryCmd.Flags().BoolVarP(&inputStdin,
		"stdin", "i",
		defaultInputStdin,
		"Read query from stdin",
	)
	queryCmd.Flags().StringSliceVarP(&textFilePaths,
		"text-files", "t",
		defaultTextFilePaths,
		"Text files added to query",
	)
	return queryCmd
}
