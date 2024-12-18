package db

import (
	"github.com/gocql/gocql"
	"github.com/sirupsen/logrus"
)

type Cassandra struct {
	Host     string
	Username string
	Password string
	Keyspace string
}

func ConnectCassandra(logger *logrus.Logger, cas *Cassandra) (*gocql.Session, error) {
	cluster := gocql.NewCluster(cas.Host)
	cluster.Keyspace = cas.Keyspace
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: cas.Username,
		Password: cas.Password,
	}
	session, err := cluster.CreateSession()
	if err != nil {
		logger.Fatalf("Failed to connect to Cassandra: %v", err)
	}
	logger.Info("Connected to Cassandra")
	return session, nil
}
