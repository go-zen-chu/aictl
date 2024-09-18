//go:generate mockgen -source=$GOFILE -destination=mock_$GOFILE -package=$GOPACKAGE
package cmd

import (
	"log/slog"
	"os"

	"github.com/go-zen-chu/aictl/usecase/query"
	"github.com/spf13/cobra"
)

<<<<<<< HEAD
var dic = di.NewContainer()

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "aictl",
	Short: "A handy cli sending query to generative AI",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) {},
}

// RootCmdExecute is function for running rootCmd (used for testing)
func RootCmdExecute() error {
	err := rootCmd.Execute()
	if err != nil {
		return fmt.Errorf("root command: %w", err)
	}
	if verbose {
		slog.Info("verbose output enabled")
		slog.SetDefault(
			slog.New(
				slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
					Level: slog.LevelDebug,
				}),
			),
		)
		slog.Debug("verbose debug output enabled")
		slog.Info("verbose info output enabled")
	}
	return nil
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := RootCmdExecute()

	if err != nil {
		slog.Error("root command failed", "error", err)
		os.Exit(1)
	}
}
=======
type CommandRequirements interface {
	UsecaseQuery() query.UsecaseQuery
}

func NewRootCmd(cmdReq CommandRequirements) *cobra.Command {
	const defaultVerbose = false
	var verbose bool

	// rootCmd represents the base command when called without any subcommands
	rootCmd := &cobra.Command{
		Use:   "aictl",
		Short: "A handy cli sending query to generative AI",
		Long: `A longer description that spans multiple lines and likely contains
	examples and usage of using your application. For example:
>>>>>>> aa56a93d30736b170f3361ac1d3f4765b9754d94

	Cobra is a CLI library for Go that empowers applications.
	This application is a tool to generate the needed files
	to quickly create a Cobra application.`,
		// Uncomment the following line if your bare application
		// has an action associated with it:
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// set logger to output to stderr because stdout is used for Generative AI response
			logHandler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
				Level: slog.LevelInfo,
			})
			if verbose {
				logHandler = slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
					Level: slog.LevelDebug,
				})
			}
			slog.SetDefault(slog.New(logHandler))
			slog.Debug("verbose debug log enabled")
		},
	}
	rootCmd.PersistentFlags().BoolVarP(
		&verbose,
		"verbose", "v",
		defaultVerbose,
		"verbose output (log level debug)",
	)
	rootCmd.AddCommand(NewQueryCmd(cmdReq))
	return rootCmd
}
