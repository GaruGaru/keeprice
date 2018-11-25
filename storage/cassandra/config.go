package cassandra

import "github.com/gocql/gocql"

type Config struct {
	Hosts             []string
	Username          string
	Password          string
	Consistency       gocql.Consistency
	ReplicationFactor int

	KeySpace string
	Table    string
}
