// Package auth handles user authentication.
package auth

import (
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/boreq/blogs/database"
	"github.com/boreq/blogs/logging"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"sync"
	"time"
)

var UsernameTakenError = errors.New("Username is already taken")
var InvalidUsernameOrPasswordError = errors.New("Invalid username or password")

var changeUsernameMutex sync.Mutex
var log = logging.New("auth")

const AuthorizationHeader = "Authorization"
const sessionKeySize = 256
const bcryptCost = 10

// GetUser returns the current user.
func GetUser(r *http.Request) (User, error) {
	// Get the session
	session := getUserSession(r)
	if session == nil {
		return &anonymousUser{}, nil
	}

	// Update LastSeen
	session.UserSession.LastSeen = time.Now()
	if _, err := database.DB.Exec(
		"UPDATE user_session SET last_seen=$1 WHERE id=$2",
		session.UserSession.LastSeen, session.UserSession.ID); err != nil {
		return nil, err
	}

	return &authenticatedUser{session.User}, nil
}

// LogoutUser logs out the current user. It is safe to call this function if
// a user is not logged in.
func LogoutUser(r *http.Request) error {
	session := getUserSession(r)
	if session != nil {
		if _, err := database.DB.Exec(
			"DELETE FROM user_session WHERE id=$1",
			session.UserSession.ID); err != nil {
			return err
		}
	}
	return nil
}

// LoginUser logs in a user. InvalidUsernameOrPasswordError is returned if the
// username or password is incorrect.
func LoginUser(username, password string) (*database.User, string, error) {
	user := getUser(username, password)
	if user == nil {
		return nil, "", InvalidUsernameOrPasswordError
	}
	sessionKey, err := createUserSession(*user)
	if err != nil {
		return nil, "", err
	}
	return user, sessionKey, nil
}

type userSession struct {
	database.User
	database.UserSession
}

func getUserSession(r *http.Request) *userSession {
	token := r.Header.Get(AuthorizationHeader)
	if token == "" {
		return nil
	}
	session := &userSession{}
	if err := database.DB.Get(session,
		`SELECT u.*, us.*
		FROM user_session us
		JOIN "user" u ON us.user_id=u.id
		WHERE us.key=$1
		LIMIT 1`,
		token); err != nil {
		return nil
	}
	return session
}

func getUser(username, password string) *database.User {
	user := &database.User{}
	if err := database.DB.
		Get(user, "SELECT * FROM \"user\" WHERE username=$1 LIMIT 1", username); err != nil {
		return nil
	}
	if !compareHashAndPassword(user.Password, password) {
		return nil
	}
	return user
}

func createUserSession(user database.User) (string, error) {
	sessionKey := make([]byte, sessionKeySize)
	if err := generateSessionKey(sessionKey); err != nil {
		return "", err
	}
	sessionKeyString := fmt.Sprintf("%x", sessionKey)
	session := database.UserSession{
		UserID:   user.ID,
		Key:      sessionKeyString,
		LastSeen: time.Now(),
	}
	if _, err := database.DB.Exec(
		`INSERT INTO user_session (user_id, key, last_seen)
		VALUES ($1, $2, $3)`,
		session.UserID, session.Key, session.LastSeen); err != nil {
		return "", err
	}
	return session.Key, nil
}

// CreateUser creates a user. UsernameTakenError is returned if the username is
// already taken.
func CreateUser(username, password string) error {
	changeUsernameMutex.Lock()
	defer changeUsernameMutex.Unlock()

	if usernameTaken(username) {
		return UsernameTakenError
	} else {
		passwordHash, err := generatePasswordHash(password)
		if err != nil {
			return err
		}
		if _, err := database.DB.Exec(
			"INSERT INTO \"user\" (username, password) VALUES ($1, $2)",
			username, passwordHash); err != nil {
			return err
		}
	}
	return nil
}

func usernameTaken(username string) bool {
	var user database.User
	err := database.DB.Get(&user, "SELECT * FROM \"user\" WHERE username=$1", username)
	return err == nil
}

func generateSessionKey(nonce []byte) error {
	_, err := rand.Read(nonce)
	return err
}

func generatePasswordHash(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	return string(b), err
}

func compareHashAndPassword(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
