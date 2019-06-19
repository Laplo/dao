package dao

import (
	"apimsprdev/models"
	"database/sql"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

func CreateUser(password string, email string) string {
	dsn := "uizcyj7nioelavez:3Q9UjGOeS6RzT6AgUgRn@tcp(bh2mwmn4sbol75r1fqss-mysql.services.clever-cloud.com)/bh2mwmn4sbol75r1fqss"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}

	stmtIns, err := db.Prepare("INSERT INTO user VALUES( ?, ?, ?, ? )") // ? = placeholder
	if err != nil {
		return ""
	}
	defer stmtIns.Close() // Close the statement when we leave main() / the program terminates
	// Next, insert the username, along with the hashed password into the database
	id, _ := uuid.NewV4()
	_, _ = stmtIns.Exec(id, email, hashedPassword, 0)

	defer db.Close()
	return id.String()
}

func GetUserByEmailAndPassword(password string, email string) models.User {
	dsn := "uizcyj7nioelavez:3Q9UjGOeS6RzT6AgUgRn@tcp(bh2mwmn4sbol75r1fqss-mysql.services.clever-cloud.com)/bh2mwmn4sbol75r1fqss"
	db, _ := sql.Open("mysql", dsn)
	user := models.User{}

	// Get the existing entry present in the database for the given username
	row := db.QueryRow("SELECT id, password, (isAdmin = b'1') FROM user WHERE email=?", email)
	var pass string
	var isAdmin bool
	var id uuid.UUID

	if err := row.Scan(&id, &pass, &isAdmin); err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("Zero rows found")
		} else {
			panic(err)
		}
		// Check for a scan error.
		// Query rows will be closed with defer.
		fmt.Println("couldn't scan id")
		return user
	}
	user.Email = email
	user.Password = pass
	user.IsAdmin = isAdmin
	user.ID = id

	// Store the obtained password in `storedCreds`
	// Compare the stored hashed password, with the hashed version of the password that was received
	if user.Password != password || len(user.Password) == 0 {
		return user
	}
	defer db.Close()
	return user
}

func DeleteUser(id uuid.UUID) {
	dsn := "uizcyj7nioelavez:3Q9UjGOeS6RzT6AgUgRn@tcp(bh2mwmn4sbol75r1fqss-mysql.services.clever-cloud.com)/bh2mwmn4sbol75r1fqss"
	db, _ := sql.Open("mysql", dsn)

	stmtDel, err := db.Prepare("DELETE FROM user WHERE id = ?")
	if err != nil {
		fmt.Println(err)
	}
	defer stmtDel.Close() // Close the statement when we leave main() / the program terminates
	// Next, insert the username, along with the hashed password into the database
	_, _ = stmtDel.Exec(id)
	defer db.Close()
}
