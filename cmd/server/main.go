package main

import (
	"fmt"
	"net/http"

	"github.com/aleroxac/goexpert-fatest-zipcode-api/configs"
	_ "github.com/aleroxac/goexpert-fatest-zipcode-api/docs"
	"github.com/aleroxac/goexpert-fatest-zipcode-api/internal/infra/webserver/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

//	@title			goexpert-fatest-zipcode-api
//	@version		1.0.0
//	@description	Segundo desafio do treinamento GoExpert(FullCycle)
//	@termsOfService	http://swagger.io/terms

//	@contact.name	Augusto Cardoso dos Santos
//	@contact.url	https://github.com/aleroxac
//	@contact.email	acardoso.ti@gmail.com

//	@license.name	Full Cycle License
//	@license.url	https://fullcycle.com.br

// @host		localhost:8080
// @BasePath	/
func main() {
	// configs
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	// HTTP Router
	r := chi.NewRouter()

	// middlewares
	r.Use(middleware.Logger)    // go-chi
	r.Use(middleware.Recoverer) // go-chi

	// cep: routes
	getCEPHandler := handlers.NewGetCEPHandler()
	r.Route("/cep", func(r chi.Router) {
		r.Get("/{cep}", getCEPHandler.GetCEP)
	})

	// swagger
	r.Get("/docs/*", httpSwagger.Handler(httpSwagger.URL("http://localhost:8080/docs/doc.json")))

	// server
	http.ListenAndServe(fmt.Sprintf(":%s", configs.WebServerPort), r)
}
