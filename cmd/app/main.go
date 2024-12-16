package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/gocql/gocql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v5"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"url-shortener/config"
	"url-shortener/internal/api"
	"url-shortener/internal/db"
	"url-shortener/pkg/models"
)
import _ "github.com/jackc/pgx/v5"

func NewPostgres(username, password, host, port, dbName string) (*pgx.Conn, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", username, password, host, port, dbName)
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func main() {
	cfg := config.Load()

	dbConn, err := db.ConnectPostgres(cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	api.StartSever(cfg.ServerPort, dbConn)

	//http.ListenAndServe(":8080", nil)

	//addUrl()
	//
	//return
	//conn, err := NewPostgres("beka", "Beka2005", "localhost", "5432", "url_shortener")
	//if err != nil {
	//	log.Fatalf("Error : %v", err)
	//}
	//newUser := models.User{
	//	Id:         1,
	//	Username:   "Bekaryss",
	//	Email:      "Beka123@gmail.com",
	//	Password:   "Beka2005",
	//	UrlCounter: 0,
	//}
	//
	//err = InsertUserIntoDatabase(conn, newUser)
	//if err != nil {
	//	return
	//}
	//
	//defer conn.Close(context.Background())
}

func getConnection() (*pgx.Conn, error) {
	DB, err := NewPostgres("beka", "Beka2005", "localhost", "5432", "url_shortener")
	if err != nil {
		return nil, err
	}
	return DB, nil
}

func addUrl() {
	cluster := gocql.NewCluster("localhost") // Replace with your Cassandra node(s)
	cluster.Keyspace = "url_shortener"       // Replace with your keyspace
	cluster.Consistency = gocql.Quorum
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: "beka",     // Your Cassandra username
		Password: "Beka2005", // Your Cassandra password
	}

	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatalf("Failed to connect to Cassandra: %v", err)
	}
	defer session.Close()

	new_url := models.Url{gocql.TimeUUID(), 1, "short", "long"}

	// Insert the record
	query := `
        INSERT INTO urls (id, user_id, short_url, long_url) 
        VALUES (?, ?, ?, ?)
    `
	err = session.Query(query, new_url.Id, new_url.UserId, new_url.ShortUrl, new_url.LongUrl).Exec()
	if err != nil {
		log.Fatalf("Failed to insert data into Cassandra: %v", err)
	}

	fmt.Println("Record added successfully!")

}

func applyMigrations(dsn string) {
	migrationDir := "migrations" // Path to migration files
	m, err := migrate.New(migrationDir, dsn)
	if err != nil {
		log.Fatalf("Failed to initialize migration: %v", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("Failed to apply migrations: %v", err)
	}
	fmt.Println("Migrations applied successfully!")
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}

//
//func connectToPostgres(dsn string) (*pgx.Conn, error) {
//
//}
