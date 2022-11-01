package database

import (
	"fmt"

	"github.com/gocql/gocql"
)

type DBconnection struct {
	Session *gocql.Session
}

var Connection DBconnection

func SetupDBconnection() {
	cluster := gocql.NewCluster("cassandra:9042")
	cluster.Keyspace = "chat"
	cluster.Consistency = gocql.Quorum
	Cs, err := cluster.CreateSession()
	Connection.Session = Cs
	if err != nil {
		fmt.Println("database error")
		panic(err)
	}
}

func ExecuteQuery(query string, args ...interface{}) error {
	err := Connection.Session.Query(query, args...).Exec() // connection.session.Close()
	return err
}

func SelectQuery(query string, args ...interface{}) *gocql.Query {
	data := Connection.Session.Query(query, args...)
	return data
}

func CheckUserExist(query string, id string) (string, string) {
	var ID string
	var Username string
	m := map[string]interface{}{}
	iter := Connection.Session.Query(query, id).Iter()
	for iter.MapScan(m) {
		ID = m["id"].(string)
		Username = m["username"].(string)
		m = map[string]interface{}{}
	}

	return ID, Username
}
