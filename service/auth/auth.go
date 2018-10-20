// Package auth handles user authentication.
package auth

import (
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/boreq/blogs/database"
	"github.com/boreq/blogs/logging"
	"github.com/boreq/sqlx"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"sync"
	"time"
)

var UsernameTakenError = errors.New("Username is already taken")
var InvalidUsernameOrPasswordError = errors.New("Invalid username or password")

var log = logging.New("auth")

const authorizationHeader = "Authorization"
const sessionKeySize = 256
const bcryptCost = 10

func New(db *sqlx.DB) *AuthService {
	rv := &AuthService{
		db: db,
	}
	return rv
}

type AuthService struct {
	db                  *sqlx.DB
	changeUsernameMutex sync.Mutex
}

// GetUser returns the current user.
func (a *AuthService) GetUser(r *http.Request) (User, error) {
	// Get the session
	session := a.getUserSession(r)
	if session == nil {
		return &anonymousUser{}, nil
	}

	// Update LastSeen
	session.UserSession.LastSeen = time.Now()
	if _, err := a.db.Exec(
		"UPDATE user_session SET last_seen=$1 WHERE id=$2",
		session.UserSession.LastSeen, session.UserSession.ID); err != nil {
		return nil, err
	}

	return &authenticatedUser{session.User}, nil
}

// LogoutUser logs out the current user. It is safe to call this function if
// a user is not logged in.
func (a *AuthService) LogoutUser(r *http.Request) error {
	session := a.getUserSession(r)
	if session != nil {
		if _, err := a.db.Exec(
			"DELETE FROM user_session WHERE id=$1",
			session.UserSession.ID); err != nil {
			return err
		}
	}
	return nil
}

// LoginUser logs in a user. InvalidUsernameOrPasswordError is returned if the
// username or password is incorrect.
func (a *AuthService) LoginUser(username, password string) (*database.User, string, error) {
	user := a.getUser(username, password)
	if user == nil {
		return nil, "", InvalidUsernameOrPasswordError
	}
	sessionKey, err := a.createUserSession(*user)
	if err != nil {
		return nil, "", err
	}
	return user, sessionKey, nil
}

type userSession struct {
	database.User
	database.UserSession
}

func (a *AuthService) getUserSession(r *http.Request) *userSession {
	token := r.Header.Get(authorizationHeader)
	if token == "" {
		return nil
	}
	session := &userSession{}
	if err := a.db.Get(session,
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

func (a *AuthService) getUser(username, password string) *database.User {
	user := &database.User{}
	if err := a.db.Get(user, "SELECT * FROM \"user\" WHERE username=$1 LIMIT 1", username); err != nil {
		return nil
	}
	if !compareHashAndPassword(user.Password, password) {
		return nil
	}
	return user
}

func (a *AuthService) createUserSession(user database.User) (string, error) {
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
	if _, err := a.db.Exec(
		`INSERT INTO user_session (user_id, key, last_seen)
		VALUES ($1, $2, $3)`,
		session.UserID, session.Key, session.LastSeen); err != nil {
		return "", err
	}
	return session.Key, nil
}

// CreateUser creates a user. UsernameTakenError is returned if the username is
// already taken.
func (a *AuthService) CreateUser(username, password string) error {
	a.changeUsernameMutex.Lock()
	defer a.changeUsernameMutex.Unlock()

	if a.usernameTaken(username) {
		return UsernameTakenError
	} else {
		passwordHash, err := generatePasswordHash(password)
		if err != nil {
			return err
		}
		if _, err := a.db.Exec(
			"INSERT INTO \"user\" (username, password) VALUES ($1, $2)",
			username, passwordHash); err != nil {
			return err
		}
	}
	return nil
}

func (a *AuthService) usernameTaken(username string) bool {
	var user database.User
	err := a.db.Get(&user, "SELECT * FROM \"user\" WHERE username=$1", username)
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
