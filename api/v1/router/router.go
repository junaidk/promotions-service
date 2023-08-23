package router

import (
	"fmt"
	middleware "github.com/deepmap/oapi-codegen/pkg/chi-middleware"
	chiMiddleware "github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"log"
	"net"
	"net/http"
	"os"
	api "promotions-service/api/v1/rest"
	"promotions-service/internal/database"
	"promotions-service/internal/processor"
	"promotions-service/internal/service"
)

func Serve(port string) {
	swagger, err := api.GetSwagger()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading swagger spec\n: %s", err)
		os.Exit(1)
	}

	swagger.Servers = nil

	repo := database.NewInMemory()
	proc := processor.NewCsvFile(repo)
	ad := service.NewAdmin(proc, repo)
	pr := service.NewPromotions(repo)

	// Create an instance of handler which satisfies the generated interface
	handler := api.NewPromotionHandler(
		&pr,
		&ad,
	)

	r := chi.NewRouter()
	r.Use(
		chiMiddleware.Recoverer,
		middleware.OapiRequestValidator(swagger),
		render.SetContentType(render.ContentTypeJSON),
	)
	api.HandlerFromMux(handler, r)

	s := &http.Server{
		Handler: r,
		Addr:    net.JoinHostPort("0.0.0.0", port),
	}

	log.Fatal(s.ListenAndServe())
}
