package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"go_with_sql/repository"
)

func printUsers(connection *sql.DB) {
	users := repository.GetUsers(connection)
	data := ToJson(users)
	fmt.Printf("#### All > %s\n", data)

	user := repository.GetUserById(connection, 4)
	data = ToJson(user)
	fmt.Printf("### ById > %s\n", data)

	user = repository.GetUserByUsername(connection, "jane")
	data = ToJson(user)
	fmt.Printf("## ByUsername > %s\n", data)

	updated := repository.UpdateUser(connection, "Jane Doe", user.Id)
	fmt.Printf("# NameUpdated > %v\n", updated)
}

// ToJson This is generics like in java
func ToJson[D any](d D) string {
	marshaled, err := json.MarshalIndent(d, "", " ")
	checkError(err)
	return string(marshaled)
}
