package main

import (
	"fmt"
	"github.com/gocql/gocql"
	_ "github.com/lib/pq"
	"log"
	"url-shortener/config"
	"url-shortener/internal/db"
	"url-shortener/internal/server"
	"url-shortener/pkg/logging"
	"url-shortener/pkg/models"
)
import _ "github.com/jackc/pgx/v5"

func main() {

	logger := logging.Init()
	cfg := config.Load(logger)

	dbConn, err := db.ConnectPostgres(cfg.DatabaseURL, logger)
	if err != nil {
		logger.Fatal(err)
	}
	cassandraSession, err := db.ConnectCassandra(logger, &cfg.Cassandra)
	if err != nil {
		logger.Fatal(err)
	}
	server.StartSever(cfg.ServerPort, dbConn, cassandraSession, logger)

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

	newUrl := models.Url{gocql.TimeUUID(), 1, "short", "long"}

	query := `
        INSERT INTO urls (id, user_id, short_url, long_url) 
        VALUES (?, ?, ?, ?)
    `
	err = session.Query(query, newUrl.Id, newUrl.UserId, newUrl.ShortUrl, newUrl.LongUrl).Exec()
	if err != nil {
		log.Fatalf("Failed to insert data into Cassandra: %v", err)
	}

	fmt.Println("Record added successfully!")

	//os.OpenFile()

}
