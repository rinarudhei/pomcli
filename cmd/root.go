/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rinarudhei/pomcli/app"
	"github.com/rinarudhei/pomcli/model"
	"github.com/rinarudhei/pomcli/session"
	"github.com/rinarudhei/pomcli/session/repository"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pomcli",
	Short: "TUI Pomodoro App",
	Long: `POMCLI is a TUI Pomodoro App built for developers or TUI enthusiasts. 
	It provides a simple interface to manage pomodoro sessions and track your productivity.
	Developers can utilize git hook to log commits or manually insert activity by using the provided cli`,
	Version: "0.1",
	RunE: func(cmd *cobra.Command, args []string) error {
		repo := repository.NewRepository()
		sqliteRepo, err := repository.NewSQLiteRepo(dbPath)
		if err != nil {
			return err
		}
		s := session.NewSession(repo, sqliteRepo, viper.GetDuration("pomodoro"), viper.GetDuration("short-break"), viper.GetDuration("long-break"))

		if len(args) >= 1 && args[0] == "activity-hook" {
			s.SqliteRepo.AddActivity(model.Activity{Message: strings.Join(args[1:], " "), CompletedAt: time.Now()})
			return nil
		}
		a, err := app.NewApp(s)
		if err != nil {
			return err
		}

		return rootAction(os.Stdout, a)
	},
}

func rootAction(out io.Writer, a *app.App) error {
	if err := a.Run(); err != nil {
		fmt.Fprintln(out, err)
		return err
	}

	return nil
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var dbPath = ""

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.pomcli.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	home, err := os.UserHomeDir()
	if err != nil {
		os.Exit(1)
	}

	// Correct path
	configDir := filepath.Join(home, ".config", "pomcli")
	dbPath = filepath.Join(configDir, "pomcli.db")

	// Create directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		os.Exit(1)
	}
	rootCmd.Flags().DurationP("pomodoro", "p", 50*time.Minute, "Pomodoro duration")
	rootCmd.Flags().DurationP("short-break", "s", 10*time.Minute, "Short break duration")
	rootCmd.Flags().DurationP("long-break", "l", 59*time.Minute, "Long break duration")

	viper.BindPFlag("pomodoro", rootCmd.Flags().Lookup("pomodoro"))
	viper.BindPFlag("short-break", rootCmd.Flags().Lookup("short-break"))
	viper.BindPFlag("long-break", rootCmd.Flags().Lookup("long-break"))
}
