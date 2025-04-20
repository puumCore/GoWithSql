package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

const (
	host     = "localhost"
	port     = 5432
	database = "syntheki"
	user     = "developer"
	password = "XUQ2As8z"
)

func ValidatePassword(hashedPwd, plainPwd string) bool {
	if hashedPwd != "" && plainPwd != "" {
		err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(plainPwd))
		CheckError(err)
		return err == nil
	} else {
		return false
	}
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func GetDbConnection() (*sql.DB, error) {
	connString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, database)

	fmt.Println("#### Obtaining connection...")
	db, err := sql.Open("postgres", connString)
	CheckError(err)

	fmt.Println("### Attempting to ping...")
	err = db.Ping()
	CheckError(err)
	fmt.Println("## Connected!")

	fmt.Println("# Selecting schema...")
	_, err = db.Exec(`SET search_path TO development`)
	CheckError(err)

	return db, err
}

func GetUserByUsername(db *sql.DB, username string) User {
	row := db.QueryRow(`SELECT customer_id, created_at, full_name, phone, email, kra_pin, username, "password" FROM default_users where username = $1`, username)
	var user User
	err := row.Scan(&user.Id, &user.CreatedAt, &user.Name, &user.Phone, &user.Email, &user.KraPin, &user.Username, &user.Password)
	if err != nil {
		return User{}
	}
	return user
}

func GetUserById(db *sql.DB, id int64) User {
	row := db.QueryRow(`SELECT customer_id, created_at, full_name, phone, email, kra_pin, username, "password" FROM default_users where customer_id = $1`, id)
	var user User
	err := row.Scan(&user.Id, &user.CreatedAt, &user.Name, &user.Phone, &user.Email, &user.KraPin, &user.Username, &user.Password)
	if err != nil {
		return User{}
	}
	return user
}

func UpdateUser(db *sql.DB, newName string, id int64) bool {
	updateResult, err := db.Exec(`UPDATE customer SET full_name = coalesce($1, null) WHERE customer_id = $2`, newName, id)
	CheckError(err)
	affected, err := updateResult.RowsAffected()
	CheckError(err)
	return affected > 0
}

func GetUsers(db *sql.DB) []User {
	row, err := db.Query(`SELECT customer_id, created_at, full_name, phone, email, kra_pin, username, "password" FROM default_users order by customer_id desc`)
	CheckError(err)

	var users []User
	for row.Next() {
		var user User
		err := row.Scan(&user.Id, &user.CreatedAt, &user.Name, &user.Phone, &user.Email, &user.KraPin, &user.Username, &user.Password)
		if err != nil {
			fmt.Println("Exception thrown in loop", err)
		} else {
			users = append(users, user)
		}
	}

	return users
}

func ToJson(users []User) string {
	marshaled, err := json.MarshalIndent(users, "", " ")
	CheckError(err)
	return string(marshaled)
}

func ToJsonAlt(u User) string {
	marshaled, err := json.MarshalIndent(u, "", " ")
	CheckError(err)
	return string(marshaled)
}

type User struct {
	Id                 int64
	CreatedAt          time.Time
	Name, Phone, Email string
	KraPin             sql.NullString
	Username, Password string
}

func CheckError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}
