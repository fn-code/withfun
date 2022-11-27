package main

import (
	"fmt"
	"log"

	"github.com/gocql/gocql"
)

type EmployeByID struct {
	ID       int
	Name     string
	Position string
}

func main() {
	// connect to the cluster
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Port = 9042
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: "cassandra",
		Password: "cassandra",
	}

	cluster.ProtoVersion = 4
	cluster.Keyspace = "test_keyspace"
	cluster.Consistency = gocql.One

	session, err := cluster.CreateSession()
	if err != nil {
		log.Printf("failed create cassandra session %v", err)
	}

	defer session.Close()

	fmt.Println("This is cassandra example")

	// insert employee by id
	if err := session.Query("INSERT INTO employee_by_id (id, name, position) VALUES (?, ?, ?)",
		"2", "Reagan", "IT").Exec(); err != nil {
		log.Println("Error : ", err)
	}

	var id int
	var name string

	if err := session.Query("SELECT id, name FROM employee_by_id WHERE id = ? LIMIT 1",
		"2").Consistency(gocql.One).Scan(&id, &name); err != nil {
		log.Println("Error : ", err)
	}
	fmt.Println("Employee By ID :", id, name)

	employee := &EmployeByID{}
	iter := session.Query("SELECT * FROM employee_by_id").Iter()
	for iter.Scan(&employee.ID, &employee.Name, &employee.Position) {
		fmt.Println(employee.ID, employee.Name, employee.Position)
	}

	if err := iter.Close(); err != nil {
		log.Println("Error : ", err)
	}

}
