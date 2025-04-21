package main

import (
	"database/sql"
	"go_with_sql/repository"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {
	connection, err := repository.GetDbConnection()
	checkError(err)
	defer func(connection *sql.DB) {
		err := connection.Close()
		repository.CheckError(err)
	}(connection)

	go printUsers(connection)
	mq()
	gRPCServer()
}

func mq() {
	_, err := getMqConnection()
	checkError(err)
}

func checkError(e error) {
	if e != nil {
		panic(e)
	}
}
