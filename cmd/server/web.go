package main

import (
  "azugo.io/azugo"
  "azugo.io/azugo/server"

  "github.com/spf13/cobra"
  "example.com/project/routes"
)

// webCmd represents the web command
var webCmd = &cobra.Command{
  Use:   "web",
  Short: "Start web server",
  Long: `Web server is the only thing you need to run,
and it takes care of all the other things for you`,
  RunE: runWeb,
}

func runWeb(cmd *cobra.Command, args []string) error {
  app, err := server.New(cmd, server.Options{
    AppName:       "Example application",
    AppVer:        Version,
  })
  if err != nil {
    return err
  }

  app.Get("/", func(ctx *azugo.Context) {
    ctx.Text("Hello, World!")
  })
  
  app.Post("/find", func(ctx *azugo.Context) {
    routes.GetUser(ctx)
  })
   
  app.Post("/check", func(ctx *azugo.Context) {
    routes.CheckUsersExistence(ctx)
  })

  app.Get("/list", func(ctx *azugo.Context) {
    routes.GetAllUsers(ctx)
  })

   app.Post("/add", func(ctx *azugo.Context) {
    routes.AddUserToTheList(ctx)
  })

  app.Put("/update", func(ctx *azugo.Context) {
    routes.PatchUser(ctx)
  })

  app.Delete("/delete", func(ctx *azugo.Context) {
    routes.DeleteUser(ctx)
  })

  server.Run(app)
  return nil
}

func init() {
  initRootCmd()
  RootCmd.AddCommand(webCmd)
}