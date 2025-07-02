// @title Azugo Training API
// @version 1.0
// @description This is a sample server.
// @host localhost:8080
// @BasePath /
package main

import (
  "fmt"
  "os"
  "go.uber.org/zap"
  "github.com/spf13/cobra"
  "example.com/project/routes"
  _ "example.com/project/docs"
)

var (
Version = "0.0.1-dev"
logger *zap.Logger
RootCmd *cobra.Command
)
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
  var err error
  logger, err = zap.NewDevelopment()
  if err != nil {
    panic(err)
  }
  routes.InitLogger(logger)
  Execute()
} 