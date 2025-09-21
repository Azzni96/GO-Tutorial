# Articles API

A simple REST API for managing articles built with Go, Gorilla Mux, and MySQL.

## Features

- Create, read, update, and delete articles
- Search articles by title or description
- Sort articles by title or creation date
- Pagination support
- MySQL database with GORM

## Prerequisites

- Go 1.25.1 or higher
- MySQL server
- Git

## Installation

1. Clone the repository:
```bash
git clone <your-repo-url>
cd GO
```

2. Install dependencies:
```bash
go mod tidy
```

3. Set up your database:
   - Create a MySQL database
   - Copy `.env.sample` to `.env` and configure your database settings:

```env
APP_PORT=8080

DB_USER=your_username
DB_PASS=your_password
DB_HOST=127.0.0.1
DB_PORT=3306
DB_NAME=your_database_name
```

## Running the Application

1. Make sure your MySQL server is running

2. Start the application:
```bash
go run .
```

The API will be available at `http://localhost:8080`

## API Endpoints

### Articles

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/` | Home page |
| GET | `/articles` | Get all articles with optional filtering |
| POST | `/articles` | Create a new article |
| GET | `/articles/{id}` | Get article by ID |
| PUT | `/articles/{id}` | Update article by ID |
| DELETE | `/articles/{id}` | Delete article by ID |

### Query Parameters for GET /articles

- `search` - Search in title or description
- `sort` - Sort by: `title`, `-title`, `createdAt`, `-createdAt`
- `page` - Page number (default: 1)
- `limit` - Items per page (default: 10, max: 100)

## Example Usage

### Create an article
```bash
curl -X POST http://localhost:8080/articles \
  -H "Content-Type: application/json" \
  -d '{
    "title": "My First Article",
    "desc": "This is a sample article",
    "content": "Full content of the article goes here..."
  }'
```

### Get all articles
```bash
curl http://localhost:8080/articles
```

### Search articles
```bash
curl "http://localhost:8080/articles?search=sample&sort=-createdAt&page=1&limit=5"
```

### Get article by ID
```bash
curl http://localhost:8080/articles/1
```

### Update an article
```bash
curl -X PUT http://localhost:8080/articles/1 \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Updated Article Title",
    "desc": "Updated description"
  }'
```

### Delete an article
```bash
curl -X DELETE http://localhost:8080/articles/1
```

## Project Structure

```
.
├── main.go          # Application entry point
├── models.go        # Data models
├── db.go           # Database connection and configuration
├── handlers.go     # HTTP request handlers
├── router.go       # Route definitions
├── utils.go        # Utility functions
├── go.mod          # Go module dependencies
├── go.sum          # Dependency checksums
├── .env            # Environment variables (not in git)
├── .env.sample     # Environment variables template
└── README.md       # This file
```

## Dependencies

- [Gorilla Mux](https://github.com/gorilla/mux) - HTTP router
- [GORM](https://gorm.io/) - ORM library for Go
- [MySQL Driver](https://github.com/go-sql-driver/mysql) - MySQL driver for Go
- [GoDotEnv](https://github.com/joho/godotenv) - Environment variable loader

## Development

To contribute to this project:

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test your changes
5. Submit a pull request

