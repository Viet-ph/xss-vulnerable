package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Viet-ph/xss-vulnerable/database"
	"github.com/Viet-ph/xss-vulnerable/handlers"
	"github.com/Viet-ph/xss-vulnerable/utils"
	"github.com/joho/godotenv"
)

type apiConfig struct {
	FileserverHits int
	JwtSecret      string
}

var ApiCfg = apiConfig{
	FileserverHits: 0,
	JwtSecret:      "",
}

func main() {
	filePathRoot := "./static"
	port := "8080"
	dbPath := "database.json"
	godotenv.Load()
	ApiCfg.JwtSecret = os.Getenv("JWT_SECRET")

	// if _, err := os.Stat(dbPath); !errors.Is(err, os.ErrNotExist) {
	// 	os.Remove(dbPath)
	// }

	//DB connection
	var err error
	database.Db, err = database.NewDB(dbPath)
	if err != nil {
		log.Printf("Error creating DB")
		return
	}

	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(filePathRoot))))
	mux.HandleFunc("POST /api/comments", handlers.CreateCommentHandler)
	mux.HandleFunc("GET /api/comments", handlers.GetAllCommentHandler)
	mux.HandleFunc("GET /api/comments/{commentID}", handlers.GetCommentById)
	mux.HandleFunc("POST /api/users", handlers.CreateUserHandler)
	mux.Handle("PUT /api/users", handlers.AuthMiddleware(http.HandlerFunc(handlers.UpdateUserHandler)))

	mux.HandleFunc("/login", handlers.UserLoginHandler)
	mux.HandleFunc("POST /logout", handlers.UserLogoutHandler)
	mux.Handle("/search", handlers.AuthMiddleware(http.HandlerFunc(handlers.SearchHandler)))
	mux.Handle("/comments", handlers.AuthMiddleware(http.HandlerFunc(handlers.CommentsHandler)))
	mux.HandleFunc("/social", handlers.SocialHandler)

	corsMux := utils.MiddlewareCors(mux)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	log.Printf("Serving files from %s on port: %s\n", filePathRoot, port)
	log.Fatal(server.ListenAndServe())
}
