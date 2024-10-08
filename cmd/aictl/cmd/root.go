//go:generate mockgen -source=$GOFILE -destination=mock_$GOFILE -package=$GOPACKAGE
package cmd

import (
	"log/slog"
	"os"

	"github.com/go-zen-chu/aictl/usecase/query"
	"github.com/spf13/cobra"
)

type CommandRequirements interface {
	UsecaseQuery() query.UsecaseQuery
}

func NewRootCmd(cmdReq CommandRequirements) *cobra.Command {
	const defaultVerbose = false
	var verbose bool

	// rootCmd represents the base command when called without any subcommands
	rootCmd := &cobra.Command{
		Use:   "aictl",
		Short: "A handy CLI for asking query to generative AI",
		Long:  `A handy CLI for asking query to generative AI`,
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
