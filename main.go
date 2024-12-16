package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gocql/gocql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"url-shortener/pkg/models"
)
import _ "github.com/jackc/pgx/v5"

//type userManager interface {
//	InsertUserIntoDatabase(user models.User) error
//	GetUserFromDatabase(id int) ( models.User, error)
//	DeleteUserFromDatabase(id int) error
//	UpdateUserFromDatabase(user  models.User) error
//	GetAllUsersFromDatabase() ([] models.User, error)
//}

func NewPostgres(username, password, host, port, dbName string) (*pgx.Conn, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", username, password, host, port, dbName)
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func hashPassword(password string) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return hashedPassword, nil
}

func main() {
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/signup", createUser)

	http.ListenAndServe(":8080", nil)

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

func createUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	conn, err := getConnection()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	user.Id, err = InsertUserIntoDatabase(conn, user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	_ = json.NewEncoder(w).Encode(&user)
	fmt.Println(user)
	w.WriteHeader(http.StatusCreated)
}

func getConnection() (*pgx.Conn, error) {
	conn, err := NewPostgres("beka", "Beka2005", "localhost", "5432", "url_shortener")
	if err != nil {
		return nil, err
	}
	return conn, nil
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

func InsertUserIntoDatabase(conn *pgx.Conn, user models.User) (int64, error) {
	// Define the SQL query to insert a new book record.
	query := `
        INSERT INTO users (username, email, password, urlCounter) VALUES (@username, @email, @password, @urlCounter) RETURNING id
    `
	// Define the named arguments for the query.
	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		return 0, err
	}
	args := pgx.NamedArgs{
		"username":   user.Username,
		"email":      user.Email,
		"password":   string(hashedPassword),
		"urlCounter": user.UrlCounter,
	}
	var id int64 = 0
	conn.QueryRow(context.Background(), query, args).Scan(&id)
	if id == 0 {
		log.Println("Error Inserting user")
		return 0, err
	}
	return id, nil
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
