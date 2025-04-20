package main

import (
	"database/sql"
	"fmt"
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

	users := repository.GetUsers(connection)
	data := repository.ToJson(users)
	fmt.Printf("All#### > %s\n", data)

	user := repository.GetUserById(connection, 4)
	data = repository.ToJsonAlt(user)
	fmt.Printf("ById### > %s\n", data)

	user = repository.GetUserByUsername(connection, "jane")
	data = repository.ToJsonAlt(user)
	fmt.Printf("ByUsername## > %s\n", data)

	updated := repository.UpdateUser(connection, "Jane Doe", user.Id)
	fmt.Printf("NameUpdated# > %v\n", updated)
}

func checkError(e error) {
	if e != nil {
		panic(e)
	}
}
