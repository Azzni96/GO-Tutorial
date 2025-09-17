package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func handleRequests() {
	r := mux.NewRouter().StrictSlash(true)

	r.HandleFunc("/", homePage).Methods("GET")
	r.HandleFunc("/articles", getAllArticles).Methods("GET")
	r.HandleFunc("/articles", createArticle).Methods("POST")
	r.HandleFunc("/articles/{id}", getArticleByID).Methods("GET")
	r.HandleFunc("/articles/{id}", updateArticle).Methods("PUT")
	r.HandleFunc("/articles/{id}", deleteArticle).Methods("DELETE")

	port := getEnv("APP_PORT", "8080")
	log.Println("🚀 Server listening on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
