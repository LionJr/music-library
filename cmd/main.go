package main

import (
	"context"
	"log"

	"github.com/chapsuk/grace"

	"music-library/internal/app"
)

// @title           Swagger Example API
// @version         1.0
// @description     Music library example.

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/

func main() {
	ctx := grace.ShutdownContext(context.Background())

	application, err := app.New(ctx)
	if err != nil {
		log.Fatalf("failed to initialize application: %v", err)
	}
	defer application.Shutdown()

	if err = application.Run(ctx); err != nil {
		log.Printf("application stopped with error: %v", err)
	}
}
