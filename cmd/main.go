package main

import (
	"log"
	"mba-golang-rate-limiter/pkg"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	rl, err := pkg.NewRateLimiter()
	if err != nil {
		log.Fatalf("Could not create rate limiter: %v", err)
	}

	r := mux.NewRouter()
	r.Use(pkg.RateLimitMiddleware(rl))

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))
	})

	http.Handle("/", r)
	log.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
