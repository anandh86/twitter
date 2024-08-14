package main

import (
	"log"
	"net/http"

	"github.com/anandh86/chirpy/internal/core/usecases"
	"github.com/anandh86/chirpy/internal/handlers"
	adapters "github.com/anandh86/chirpy/internal/repositories"
	"github.com/go-chi/chi"
)

func main() {

	// wiring
	userRepository := adapters.ProvideInMemoryRepo()
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
	subRouter.Put("/users", userHttpHandler.UpdateUser)
	subRouter.Post("/login", userHttpHandler.LoginUser)
	subRouter.Post("/refresh", userHttpHandler.Refresh)
	subRouter.Post("/revoke", userHttpHandler.Revoke)

	subRouter.Post("/tweets", userHttpHandler.PostTweet)
	subRouter.Get("/tweets/{tweetId}", userHttpHandler.GetTweetById)
	subRouter.Get("/tweets", userHttpHandler.GetAllTweets)
	subRouter.Delete("/tweets/{tweetId}", userHttpHandler.DeleteTweet)

	subRouter.Post("/polka/webhooks", userHttpHandler.PolkaWebHooks)

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
