package main

import (
	"fmt"

	"github.com/nikumar1206/puff"
	"github.com/nikumar1206/puff/middleware"
	"github.com/tiredkangaroo/sculpt"
)

type AssetsInput struct {
	FileName string `kind:"path"`
}

func startServer(mode string, item *sculpt.Model) {
	app := puff.DefaultApp("Checklist")

	fFR := puff.FileResponse{FilePath: "static/index.html"}
	app.Get("/", "Serves static application.", nil, fFR.Handler())

	assetsInput := new(AssetsInput)
	app.Get("/assets/{file}", "Serves assets for static application.", assetsInput, func(c *puff.Context) {
		// FIXME: potentially dangerous
		c.SendResponse(puff.FileResponse{FilePath: "static/assets/" + assetsInput.FileName})
	})

	apiRouter := getAPIRouter(item)

	if mode == "development" {
		fmt.Println("Running server in debug mode 🚀")
		app.Use(middleware.CORS())
	}

	app.IncludeRouter(apiRouter)
	app.ListenAndServe(":8030")
}
