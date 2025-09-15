package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type Article struct {
	Id      string `json:"Id"`
	Title   string `json:"title"`
	Desc    string `json:"desc"`
	Content string `json:"content"`
}

var articles = []Article{
	{Id: "1", Title: "Hello", Desc: "Article Description", Content: "Article Content"},
	{Id: "2", Title: "Hello 2", Desc: "Article Description", Content: "Article Content"},
}

func writeJSON(
	w http.ResponseWriter,
	status int,
	v any,
) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
func writeError(
	w http.ResponseWriter,
	status int,
	msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func newID() string {
	return strconv.FormatInt(time.Now().UnixNano(), 10)
}

func homePage(
	w http.ResponseWriter,
	r *http.Request,
) {
	fmt.Fprint(w, "Welcome to the HomePage!")
}
func getAllArticles(
	w http.ResponseWriter,
	r *http.Request,
) {
	writeJSON(
		w,
		http.StatusOK,
		articles,
	)
}
func getArticleByID(
	w http.ResponseWriter,
	r *http.Request,
) {
	id := mux.Vars(r)["id"]
	for _, a := range articles {
		if a.Id == id {
			writeJSON(
				w,
				http.StatusOK,
				a,
			)
			return
		}
	}
	writeError(
		w,
		http.StatusNotFound,
		"Article not found",
	)
}

func createArticle(
	w http.ResponseWriter,
	r *http.Request,
) {
	body, err := io.ReadAll(
		r.Body,
	)
	if err != nil {
		writeError(
			w,
			http.StatusBadRequest,
			"Invalid request body",
		)
		return
	}
	var a Article
	if err :=
		json.Unmarshal(
			body,
			&a,
		); err != nil {
		writeError(
			w,
			http.StatusBadRequest,
			"Invalid JSON",
		)
		return
	}
	if a.Id == "" {
		a.Id = newID()
	}
	articles = append(articles, a)
	writeJSON(
		w,
		http.StatusCreated,
		a,
	)
}
func updateArticle(
	w http.ResponseWriter,
	r *http.Request,
) {
	id := mux.Vars(r)["id"]
	body, err := io.ReadAll(
		r.Body,
	)
	if err != nil {
		writeError(
			w,
			http.StatusBadRequest,
			"Invalid request body",
		)
		return
	}
	var payload Article
	if err :=
		json.Unmarshal(
			body,
			&payload,
		); err != nil {
		writeError(
			w,
			http.StatusBadRequest,
			"Invalid JSON",
		)
		return
	}
	for i, a := range articles {
		if a.Id == id {
			if payload.Title != "" {
				a.Title = payload.Title
			}
			if payload.Desc != "" {
				a.Desc = payload.Desc
			}
			if payload.Content != "" {
				a.Content = payload.Content
			}
			a.Id = id
			articles[i] = a
			writeJSON(
				w,
				http.StatusOK,
				a,
			)
			return
		}
	}
	writeError(
		w,
		http.StatusNotFound,
		"Article not found",
	)
}
func deleteArticle(
	w http.ResponseWriter,
	r *http.Request,
) {
	id := mux.Vars(r)["id"]
	for i, a := range articles {
		if a.Id == id {
			articles = append(articles[:i], articles[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	writeError(
		w,
		http.StatusNotFound,
		"Article not found",
	)
}
func handleRequests() {
	r := mux.NewRouter().StrictSlash(true)

	r.HandleFunc("/", homePage).Methods("GET")

	r.HandleFunc("/articles", getAllArticles).Methods("GET")
	r.HandleFunc("/articles", createArticle).Methods("POST")
	r.HandleFunc("/articles/{id}", getArticleByID).Methods("GET")
	r.HandleFunc("/articles/{id}", updateArticle).Methods("PUT")
	r.HandleFunc("/articles/{id}", deleteArticle).Methods("DELETE")

	log.Println("Server listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
func main() {
	handleRequests()
}
