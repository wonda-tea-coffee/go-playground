package main

import (
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

func CustomHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Custom-Header", "hai")

		next.ServeHTTP(w, r)
	})
}

func SimpleHandler(w http.ResponseWriter, _ *http.Request) {
	_, err := io.WriteString(w, "OK")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "INTERNAL SERVER ERROR", 500)
	}
}

func setupRouter() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/", SimpleHandler).Methods(http.MethodGet)

	composedCommonMiddleware := alice.New(
		CustomHandler,
	)
	r.Use(composedCommonMiddleware.Then)

	return r
}

func Run() error {
	srv := &http.Server{
		Addr:              ":8080",
		Handler:           setupRouter(),
		ReadHeaderTimeout: 30 * time.Second,
	}
	return srv.ListenAndServe()
}

func main() {
	if err := Run(); err != nil {
		log.Fatal(err)
	}
}
