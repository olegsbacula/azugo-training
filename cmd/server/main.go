package main

import (
  "fmt"
  "os"

  "github.com/spf13/cobra"
)

// Version holds the current application version.
//
// This can be set using build tag to set real version number.
var Version = "0.0.1-dev"

// RootCmd represents the base command when called without any subcommands
var RootCmd *cobra.Command

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
  if err := RootCmd.Execute(); err != nil {
    fmt.Println(err)
    os.Exit(-1)
  }
}

func initRootCmd() {
  if RootCmd != nil {
    return
  }
  RootCmd = &cobra.Command{
    Use:   "server",
    Short: "Server",
    Long: `By default, server will start serving using the web server with no
  arguments - which can alternatively be run by running the subcommand web.`,
    RunE: runWeb,
  }
}

func main() {
  initRootCmd() // to run do: "go run ./cmd/server"
  RootCmd.Version = Version
  Execute()
}