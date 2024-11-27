package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
	"log"
)
import _ "github.com/jackc/pgx/v5"

type user struct {
	id         int
	username   string
	email      string
	password   string
	urlCounter int
}

type userManager interface {
	InsertUserIntoDatabase(user user) error
	GetUserFromDatabase(id int) (user, error)
	DeleteUserFromDatabase(id int) error
	UpdateUserFromDatabase(user user) error
	GetAllUsersFromDatabase() ([]user, error)
}

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
	conn, err := NewPostgres("beka", "Beka2005", "localhost", "5432", "url_shortener")
	if err != nil {
		log.Fatalf("Error : %v", err)
	}
	newUser := user{
		id:         1,
		username:   "Bekaryss",
		email:      "Beka123@gmail.com",
		password:   "Beka2005",
		urlCounter: 0,
	}
	err = InsertUserIntoDatabase(conn, newUser)
	if err != nil {
		return
	}

	defer conn.Close(context.Background())
}

func InsertUserIntoDatabase(conn *pgx.Conn, user user) error {
	// Define the SQL query to insert a new book record.
	query := `
        INSERT INTO users (username, email, password, urlCounter) VALUES (@username, @email, @password, @urlCounter)
    `
	// Define the named arguments for the query.
	hashedPassword, err := hashPassword(user.password)
	if err != nil {
		return err
	}
	args := pgx.NamedArgs{
		"username":   user.username,
		"email":      user.email,
		"password":   string(hashedPassword),
		"urlCounter": user.urlCounter,
	}
	_, err = conn.Exec(context.Background(), query, args)
	if err != nil {
		log.Println("Error Inserting user")
		fmt.Println(err)
		return err
	}
	return nil
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
