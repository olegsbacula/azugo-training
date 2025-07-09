package main

import (
	"azugo.io/azugo"
	"azugo.io/azugo/server"
	"example.com/project/routes"
	"github.com/spf13/cobra"
	"github.com/valyala/fasthttp"


	"github.com/prometheus/client_golang/prometheus"

)

// webCmd represents the web command
var (
	webCmd = &cobra.Command{
		Use:   "web",
		Short: "Start web server",
		Long: `Web server is the only thing you need to run,
and it takes care of all the other things for you`,
		RunE: runWeb,
	}
)

var (
	opsProcessed = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "myapp_processed_ops_total",
			Help: "The total number of processed operations",
		},
	)
)

// func corsMiddleware(next azugo.RequestHandler) azugo.RequestHandler {
// 	return func(ctx *azugo.Context) {

// 		ctx.Header.Set("Access-Control-Allow-Origin", "*")
// 		ctx.Header.Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
// 		ctx.Header.Set("Access-Control-Allow-Headers", "Content-Type,Authorization")

// 		if string(ctx.Method()) == "OPTIONS" {
// 			ctx.StatusCode(200)
// 			return
// 		}
// 		next(ctx)
// 	}
// }

func runWeb(cmd *cobra.Command, args []string) error {
	fsDocs := &fasthttp.FS{
		Root:               "./docs",
		IndexNames:         []string{"swagger.json"},
		GenerateIndexPages: false,
	}
	docsHandler := fsDocs.NewRequestHandler()
	app, err := server.New(cmd, server.Options{
		AppName: "Example application",
		AppVer:  Version,
	})
	if err != nil {
		return err
	}
	//app.Use(corsMiddleware)

	app.Post("/find/{username}", func(ctx *azugo.Context) {
    routes.GetUser(ctx)
})

	app.Post("/check", func(ctx *azugo.Context) {
		routes.CheckUsersExistence(ctx)
		opsProcessed.Inc()
	})

	app.Get("/list", func(ctx *azugo.Context) {
		routes.GetAllUsers(ctx)
		opsProcessed.Inc()
	})

	app.Post("/add", func(ctx *azugo.Context) {
		routes.AddUserToTheList(ctx)
		opsProcessed.Inc()
	})

	app.Put("/update", func(ctx *azugo.Context) {
		routes.PatchUser(ctx)
		opsProcessed.Inc()
	})

	app.Delete("/delete", func(ctx *azugo.Context) {
		routes.DeleteUser(ctx)
		opsProcessed.Inc()
	})

	app.Get("/swagger.json", func(ctx *azugo.Context) {
		ctx.Context().URI().SetPath("swagger.json")
		docsHandler(ctx.Context())
	})
	  
	app.Get("/health", func(ctx *azugo.Context) {
        ctx.StatusCode(fasthttp.StatusOK)
        ctx.ContentType("application/json")
        ctx.JSON(map[string]string{"status": "ok"})
    })

	app.Get("/readyz", func(ctx *azugo.Context) {
		opsProcessed.Inc()
		ctx.StatusCode(fasthttp.StatusOK)
		ctx.ContentType("text/plain")
		ctx.Text("All is working")
	})

	  app.Get("/login/{id}/{password}", func(ctx *azugo.Context) {
       routes.LoginByID(ctx)
       opsProcessed.Inc()
    })

	fsUI := &fasthttp.FS{
		Root:               "./swagger",
		IndexNames:         []string{"index.html"},
		GenerateIndexPages: false,
	}
	uiHandler := fsUI.NewRequestHandler()

	app.Get("/swagger", func(ctx *azugo.Context) {
		ctx.Redirect("/swagger/index.html")
	})

	app.Get("/swagger/{filepath:*}", func(ctx *azugo.Context) {
		ctx.Context().URI().SetPath(ctx.UserValue("filepath").(string))
		uiHandler(ctx.Context())
	})

	server.Run(app)
	return nil
}

func init() {
	prometheus.MustRegister(opsProcessed)
	initRootCmd()
	RootCmd.AddCommand(webCmd)
}
