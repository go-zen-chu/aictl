package cmd

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/user"

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

			if err := printResult(res); err != nil {
				return fmt.Errorf("print result: %w", err)
			}
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
		slog.Info("Running in GitHub Actions detected")

		slog.Debug("Checking permissions")
		currentUser, err := user.Current()
		if err != nil {
			return fmt.Errorf("get current user: %w", err)
		}
		slog.Debug("Current process user name and id", "user", currentUser.Username, "uid", currentUser.Uid)

		// you will get a value like /home/runner/work/_temp/_runner_file_commands/set_output_<guid>
		ghOutputFilePath := os.Getenv("GITHUB_OUTPUT")
		if ghOutputFilePath == "" {
			return fmt.Errorf("GITHUB_OUTPUT environment variable is not set")
		}
		fi, err := os.Stat(ghOutputFilePath)
		if err != nil {
			return fmt.Errorf("GITHUB_OUTPUT file (%s) does not exist: %w", ghOutputFilePath, err)
		}
		slog.Debug("file info", "fi", fi.Sys())

		f, err := os.OpenFile(ghOutputFilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			return fmt.Errorf("open file of GITHUB_OUTPUT: %w", err)
		}
		defer f.Close()

		const githubActionOutputsResponse = "response"
		_, err = fmt.Fprintf(f, "%s<<AICTL_EOF\n%s\nAICTL_EOF\n",
			githubActionOutputsResponse,
			res,
		)
		if err != nil {
			return fmt.Errorf("write response to file of GITHUB_OUTPUT: %w", err)
		}
		slog.Debug("Write to GITHUB_OUTPUT file", "filepath", ghOutputFilePath, "content", githubActionOutputsResponse)
		bt, err := os.ReadFile(ghOutputFilePath)
		if err != nil {
			return fmt.Errorf("read github output file: %w", err)
		}
		fmt.Printf("Write query output to ${{ steps.<step_id>.outputs.response }}\nfilepath: %s\n%s\n", ghOutputFilePath, string(bt))
		return nil
	}
	// Print the response to **stdout**. All the logs are printed to **stderr**
	fmt.Println(res)
	return nil
}
