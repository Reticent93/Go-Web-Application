package main

import (

	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)
func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		req := fmt.Sprintf("%s %s", r.Method, r.URL)
		log.Println(req)
		next.ServeHTTP(w, r)
		log.Println(req, "completed in", time.Now().Sub(start))
	})
}

var templates = template.Must(template.ParseFiles("./templates/base.html", "./templates/body.html"))

func index() http.Handler {
	return  http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b := struct {
			Title template.HTML
			BusinessName string
			Slogan string
		}{
			Title: template.HTML("Business &verbar; Landing"),
			BusinessName: "Business,",
			Slogan: "We get things done",
		}
		err := templates.ExecuteTemplate(w, "base", &b)
		if err != nil {
			http.Error(w, fmt.Sprintf("index: could not parse template: %v", err), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
}

func public() http.Handler  {
	return http.StripPrefix("/public/", http.FileServer(http.Dir("./public")))
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/public/", logging(public()))
	mux.Handle("/", logging(index()))

	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8080"
	}
	addr := fmt.Sprintf(":%s", port)
	server := http.Server{
		Addr: addr,
		Handler: mux,
		ReadHeaderTimeout: 15 *time.Second,
		WriteTimeout: 15 *time.Second,
		IdleTimeout:  15 *time.Second,
	}

	log.Println("main: running on simple port", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("main: could not start simple server: %v\n", err)
	}
}
