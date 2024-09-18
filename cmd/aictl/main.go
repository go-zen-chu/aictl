package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/go-zen-chu/aictl/cmd/aictl/cmd"
	"github.com/go-zen-chu/aictl/internal/di"
	"github.com/spf13/cobra"
)

func main() {
	app := NewApp()
	if err := app.Run(os.Args); err != nil {
		slog.Error("failed while running app", "error", err)
		os.Exit(1)
	}
}

type app struct {
	rootCmd *cobra.Command
}

func NewApp() *app {
	return &app{
		rootCmd: cmd.NewRootCmd(di.NewContainer()),
	}
}

func (a *app) Run(args []string) error {
	a.rootCmd.SetArgs(args[1:])
	if err := a.rootCmd.Execute(); err != nil {
		return fmt.Errorf("root command: %w", err)
	}
	return nil
}
