package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// ====== DB ======
var db *gorm.DB

// تحميل ملف .env
func loadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ لم يتم العثور على ملف .env، سيتم استخدام متغيرات البيئة الحالية")
	}
}

func connectDB() {
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	name := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&loc=Local",
		user, pass, host, port, name)

	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("DB connect error: %v", err)
	}

	if err := db.AutoMigrate(&Article{}); err != nil {
		log.Fatalf("AutoMigrate error: %v", err)
	}

	log.Println("✅ Database connected successfully")
}

// ====== MODELS ======
type Article struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Title     string    `json:"title"  gorm:"type:varchar(200);not null"`
	Desc      string    `json:"desc"   gorm:"type:varchar(500)"`
	Content   string    `json:"content" gorm:"type:text"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// ====== HELPERS ======
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func parseIntDefault(s string, def int) int {
	if s == "" {
		return def
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return i
}

// ====== HANDLERS ======
func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "API OK. Try: GET /articles")
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

// ====== ROUTER + SERVER ======
func handleRequests() {
	r := mux.NewRouter().StrictSlash(true)

	r.HandleFunc("/", homePage).Methods("GET")
	r.HandleFunc("/articles", getAllArticles).Methods("GET")
	r.HandleFunc("/articles", createArticle).Methods("POST")
	r.HandleFunc("/articles/{id}", getArticleByID).Methods("GET")
	r.HandleFunc("/articles/{id}", updateArticle).Methods("PUT")
	r.HandleFunc("/articles/{id}", deleteArticle).Methods("DELETE")

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("🚀 Server listening on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func main() {
	loadEnv()
	connectDB()
	handleRequests()
}
