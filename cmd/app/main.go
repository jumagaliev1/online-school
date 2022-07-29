package main

import (
	"html/template"
	"net/http"
	"online-school/main/internal/database"
	"online-school/main/internal/handler"
	"online-school/main/internal/middleware"
	"online-school/main/internal/server"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {
	db := database.Connect()
	handlers := &handler.Handler{
		DB:   db,
		Tmpl: template.Must(template.ParseGlob("./web/templates/*")),
	}

	r := mux.NewRouter()

	headPage := http.NewServeMux()
	adminMux := http.NewServeMux()
	headPage.HandleFunc("/", handlers.Head)
	adminMux.HandleFunc("/admin", handlers.AdminPage)
	adminMux.HandleFunc("/addAdmin", handlers.AddAdmin)
	//set middleware for auth
	adminMiddleware := middleware.AdminAuthMiddleware(adminMux)
	authMiddleware := middleware.AuthMiddleware(headPage)
	r.HandleFunc("/login", handlers.Login).Methods("GET")
	r.HandleFunc("/auth", handlers.Auth).Methods("POST")
	r.Handle("/", authMiddleware).Methods("GET")
	r.HandleFunc("/logout", handlers.Logout).Methods("GET")
	r.HandleFunc("/addAdmin", handlers.AddAdmin).Methods("POST")
	r.Handle("/admin", adminMiddleware).Methods("GET")
	//r.Handle("/web/", http.StripPrefix("/web/", http.FileServer(http.Dir("./web"))))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./web"))))
	//http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./web/css"))))

	port := "8081"
	server.ServeHTTP(r, port)

}
