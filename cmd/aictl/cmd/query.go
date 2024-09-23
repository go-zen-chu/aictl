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

			printResult(res)
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
		"An array of text files added to query seperated with comma (e.g. file1.go,file2.txt)",
	)
	return queryCmd
}

func printResult(res string) error {
	if githubAction := os.Getenv("GITHUB_ACTIONS"); githubAction != "" && githubAction == "true" {
		outputFile := os.Getenv("GITHUB_OUTPUT")
		if outputFile == "" {
			return fmt.Errorf("GITHUB_OUTPUT environment variable is not set")
		}
		f, err := os.OpenFile(outputFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			return fmt.Errorf("open GITHUB_OUTPUT env var file: %w", err)
		}
		defer f.Close()

		const githubActionOutputsResponse = "response"
		_, err = fmt.Fprintf(f, "%s<<AICTL_EOF\n%s\nAICTL_EOF\n",
			githubActionOutputsResponse,
			res,
		)
		if err != nil {
			return fmt.Errorf("write response to GITHUB_OUTPUT env var file: %w", err)
		}
		fmt.Println("Print output to ${{ steps.<step_id>.outputs.response }}")
		return nil
	}
	// Print the response to **stdout**. All the logs are printed to **stderr**
	fmt.Println(res)
	return nil
}
