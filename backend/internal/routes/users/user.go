package users

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aaaxpel/album/internal/db"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/go-chi/jwtauth/v5"
)

type User struct {
	id         int
	username   string
	password   string
	role       string
	created_at time.Time
}

func Register(w http.ResponseWriter, r *http.Request) {

	pool := db.Connect()

	username := r.PostFormValue("username")
	password := r.PostFormValue("password")

	// validate password length
	// validate username is unique (select)

	var exists bool
	err := pool.QueryRow(
		context.Background(),
		"SELECT EXISTS(SELECT 1 FROM users WHERE username = $1);",
		username,
	).Scan(&exists)

	if err != nil {
		log.Printf("Error on existance check: %v", err)
	}

	if exists {
		// User exists
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to hash the password: %v\n", err)
	}

	_, err = pool.Exec(context.Background(), "INSERT INTO users (username, password) VALUES ('"+username+"', '"+string(hash)+"');")
	if err != nil {
		log.Printf("Error while registering: \n%v", err)
	}
}

func Login(w http.ResponseWriter, r *http.Request) {

	// Check if the user is logged in first

	username := r.PostFormValue("username")
	password := r.PostFormValue("password")

	pool := db.Connect()

	user := User{}
	err := pool.QueryRow(
		context.Background(),
		"SELECT id, username, password, role, created_at FROM users WHERE username = $1",
		username,
	).Scan(&user.id, &user.username, &user.password, &user.role, &user.created_at)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// User not found
			println("user not found")
		}
	}

	// Probably need to convert responses to JSON
	// Will consult with frontend about this

	err = bcrypt.CompareHashAndPassword([]byte(user.password), []byte(password))
	if err != nil {
		// Incorrect password
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}

	var tokenAuth *jwtauth.JWTAuth

	if secret, ok := os.LookupEnv("DB"); ok {
		tokenAuth = jwtauth.New("HS256", []byte(secret), nil)
		_, tokenString, _ := tokenAuth.Encode(map[string]interface{}{"id": &user.id, "username": &user.username, "role": &user.id})
		fmt.Println(tokenString)
	}

}

// func validate() {

// }
