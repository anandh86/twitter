package main

import (
	"log"
	"net/http"

	"github.com/anandhmaps/chirpy/internal/core/usecases"
	"github.com/anandhmaps/chirpy/internal/handlers"
	adapters "github.com/anandhmaps/chirpy/internal/repositories"
	"github.com/go-chi/chi"
)

func main() {

	// wiring
	userRepository := adapters.ProvideUserRepository()
	userUseCase := usecases.ProvideUserUseCase(userRepository)
	userHttpHandler := handlers.ProvideUserHttpHandler(userUseCase)

	const filepathRoot = "."
	const port = "8080" // Set your desired port

	r := chi.NewRouter()

	// Setup CORS middleware
	r.Use(middlewareCors)

	// Define routes
	subRouter := chi.NewRouter()

	subRouter.Post("/users", userHttpHandler.CreateUser)

	r.Mount("/api", subRouter)

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(http.ListenAndServe(":"+port, r))

}

func middlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
