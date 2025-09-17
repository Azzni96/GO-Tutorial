package main

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("API OK. Try: GET /articles"))
}

// GET /articles?search=&sort=&page=&limit=
func getAllArticles(w http.ResponseWriter, r *http.Request) {
	search := strings.TrimSpace(r.URL.Query().Get("search"))
	sortQ := strings.TrimSpace(r.URL.Query().Get("sort"))
	page := parseIntDefault(r.URL.Query().Get("page"), 1)
	limit := parseIntDefault(r.URL.Query().Get("limit"), 10)

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	q := db.Model(&Article{})

	// Filtering
	if search != "" {
		like := "%" + search + "%"
		q = q.Where("title LIKE ? OR `desc` LIKE ?", like, like)
	}

	// Sorting
	switch sortQ {
	case "title":
		q = q.Order("title asc")
	case "-title":
		q = q.Order("title desc")
	case "createdAt":
		q = q.Order("created_at asc")
	case "-createdAt":
		q = q.Order("created_at desc")
	default:
		q = q.Order("id asc")
	}

	offset := (page - 1) * limit
	var items []Article
	if err := q.Offset(offset).Limit(limit).Find(&items).Error; err != nil {
		writeError(w, http.StatusInternalServerError, "db error")
		return
	}
	writeJSON(w, http.StatusOK, items)
}

// GET /articles/{id}
func getArticleByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var a Article
	if err := db.First(&a, id).Error; err != nil {
		writeError(w, http.StatusNotFound, "Article not found")
		return
	}
	writeJSON(w, http.StatusOK, a)
}

// POST /articles
func createArticle(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	var a Article
	if err := json.Unmarshal(body, &a); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}
	if strings.TrimSpace(a.Title) == "" {
		writeError(w, http.StatusBadRequest, "title is required")
		return
	}
	if err := db.Create(&a).Error; err != nil {
		writeError(w, http.StatusInternalServerError, "db error")
		return
	}
	writeJSON(w, http.StatusCreated, a)
}

// PUT /articles/{id}
func updateArticle(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	var payload Article
	if err := json.Unmarshal(body, &payload); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	var a Article
	if err := db.First(&a, id).Error; err != nil {
		writeError(w, http.StatusNotFound, "Article not found")
		return
	}

	if strings.TrimSpace(payload.Title) != "" {
		a.Title = payload.Title
	}
	if strings.TrimSpace(payload.Desc) != "" {
		a.Desc = payload.Desc
	}
	if strings.TrimSpace(payload.Content) != "" {
		a.Content = payload.Content
	}

	if err := db.Save(&a).Error; err != nil {
		writeError(w, http.StatusInternalServerError, "db error")
		return
	}
	writeJSON(w, http.StatusOK, a)
}

// DELETE /articles/{id}
func deleteArticle(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if err := db.Delete(&Article{}, id).Error; err != nil {
		writeError(w, http.StatusInternalServerError, "db error")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
