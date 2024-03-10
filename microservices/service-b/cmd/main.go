package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/christianferraz/tempwithopentelemetry/microservice/service-b/internal/handler"
	"github.com/christianferraz/tempwithopentelemetry/microservice/service-b/pkg/otel"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	_, err := otel.Initialize("http://zipkin:9411/api/v2/spans", "service-a")
	if err != nil {
		log.Fatal(err)
	}
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/cep/{cep}", handler.BuscaCepHandler)

	fmt.Println("Starting web server on port: 8081")

	if err := http.ListenAndServe(":8081", r); err != nil {
		fmt.Println(err)
	}
}
