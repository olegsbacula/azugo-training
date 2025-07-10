// @title Swagger Example API
// @version 1.0
// @description This is a simple GO API.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @BasePath /
package main

import (
  "fmt"
  "os"
  "github.com/spf13/cobra"
  _ "example.com/project/docs"
)

var (
Version = "0.0.1-dev"
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
  if err != nil {
    panic(err)
  }
  Execute()
} 