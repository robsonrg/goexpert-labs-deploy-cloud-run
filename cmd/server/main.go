package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/robsonrg/goexpert-labs-deploy-cloud-run/internal/infra/webserver/handlers"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/joho/godotenv"
)

func init() {
	// Tenta carregar as variáveis do arquivo .env
	if err := godotenv.Load(); err != nil {
		// Não é um erro fatal se o arquivo não existir
		// pois as variáveis podem estar no ambiente real
		fmt.Println("Não foi possível carregar o arquivo .env, buscando variáveis do sistema.")
	}
}

func main() {
	_, ok := os.LookupEnv("WEATHER_API_KEY")
	if !ok {
		panic("WEATHER_API_KEY not set")
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Route("/cep-temperature", func(r chi.Router) {
		r.Method(http.MethodGet, "/{cep}", handlers.GetCepTemperatureHandler())
	})
	http.ListenAndServe(":8080", r)
}
