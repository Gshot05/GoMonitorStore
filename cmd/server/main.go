// cmd/server/main.go

package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"myapp/internal/auth"
	"myapp/internal/handler"
	"myapp/internal/service"
	"myapp/migrations"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/nats-io/nats.go"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	migrationsPath := "migrations"
	err = migrations.RunMigrations(db, migrationsPath)
	if err != nil {
		log.Fatalf("Error running migrations: %v", err)
	}

	natsURL := os.Getenv("NATS_URL")
	log.Println("Подключение к NATS...")
	nc, err := nats.Connect(natsURL)
	if err != nil {
		log.Fatalf("Error connecting to NATS: %v", err)
	} else {
		log.Print("Подключение успешно!")
	}
	defer nc.Close()

	authService := auth.NewAuthService(db, nc)
	displayService := service.NewDisplayService(db, nc)
	monitorService := service.NewMonitorService(db, nc)

	h := handler.NewHandler(
		authService,
		displayService,
		monitorService)

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/register.html")
	})
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/login.html")
	})
	http.HandleFunc("/main", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/main.html")
	})

	http.HandleFunc("/api/register", h.Register)
	http.HandleFunc("/api/login", h.Login)
	http.HandleFunc("/api/displays", h.AddDisplay)
	http.HandleFunc("/api/monitors", h.AddMonitor)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server is running on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
