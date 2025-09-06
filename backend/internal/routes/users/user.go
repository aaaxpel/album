package users

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/aaaxpel/album/internal/db"
	"golang.org/x/crypto/bcrypt"
)

func Register(w http.ResponseWriter, r *http.Request) {

	username := r.PostFormValue("username")
	password := r.PostFormValue("password")

	// validate password length
	// validate username is unique (select)

	hash, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to hash the password: %v\n", err)
	}

	println(username, string(hash))

	pool := db.Connect()

	_, err = pool.Exec(context.Background(), "INSERT INTO users (username, password) VALUES ('"+username+"', '"+string(hash)+"');")
	fmt.Println("INSERT INTO users (username, password) VALUES ('" + username + "', '" + string(hash) + "');")
	fmt.Println(err)
}

// func validate() {

// }
