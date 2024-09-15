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
	RunE: runE,
}

type inputQuery struct {
	query            string
	outputFormat     string
	responseLanguage string
	inputStdin       bool
	textFilePaths    []string
	useGitDiff       bool
}

var defaultInput = inputQuery{
	query:            "this is a test query",
	outputFormat:     "text",
	responseLanguage: "English",
	inputStdin:       false,
	textFilePaths:    []string{},
	useGitDiff:       false,
}

var input inputQuery

func init() {
	rootCmd.AddCommand(queryCmd)
	input = defaultInput

	queryCmd.Flags().StringVarP(&input.outputFormat,
		"output", "o",
		defaultInput.outputFormat,
		"Output format text or json (default is text)",
	)
	queryCmd.Flags().StringVarP(&input.responseLanguage,
		"language", "l",
		defaultInput.responseLanguage,
		"Which language you want to get response (default is English)",
	)
	queryCmd.Flags().BoolVarP(&input.inputStdin,
		"stdin", "i",
		defaultInput.inputStdin,
		"Read query from stdin",
	)
	queryCmd.Flags().StringSliceVarP(&input.textFilePaths,
		"text-files", "t",
		defaultInput.textFilePaths,
		"Text files added to query",
	)
}

func validate(args []string) error {
	if !input.inputStdin && len(args) != 1 {
		return fmt.Errorf("query command requires only 1 argument `query text`")
	}
	if input.inputStdin && len(args) != 0 {
		return fmt.Errorf("query command with stdin cannot have a `query text`")
	}
	if input.outputFormat == "" {
		return fmt.Errorf("output format is required but got empty")
	}
	of := input.outputFormat
	if of != "text" && of != "json" {
		return fmt.Errorf("output format must be text or json but got: %s", of)
	}
	return nil
}

func runE(cmd *cobra.Command, args []string) error {
	if err := validate(args); err != nil {
		return fmt.Errorf("validation in query: %w", err)
	}
	uq := dic.UsecaseQuery()
	if input.inputStdin {
		stdin, err := io.ReadAll(os.Stdin)
		if err != nil {
			return fmt.Errorf("read stdin: %w", err)
		}
		input.query = string(stdin)
	} else {
		input.query = args[0]
	}
	res, err := uq.QueryToOpenAI(
		input.query,
		input.outputFormat,
		input.responseLanguage,
		input.textFilePaths,
	)
	if err != nil {
		return fmt.Errorf("query to openai: %w", err)
	}
	// Print the response to **stdout**. All the logs are printed to **stderr**
	fmt.Printf("%s", res)
	return nil
}
