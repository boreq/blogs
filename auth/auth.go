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
var log = logging.GetLogger("auth")

const sessionKeyCookieName = "session_key"
const sessionKeySize = 256
const bcryptCost = 10

type User interface {
	// IsAuthenticated returns true if the user has logged in.
	IsAuthenticated() bool

	// GetUser returns an associated database entry or nil if the user is
	// not authenticated.
	GetUser() *database.User
}

// authenticatedUser is used to represent a logged in user.
type authenticatedUser struct {
	user database.User
}

func (u authenticatedUser) IsAuthenticated() bool {
	return true
}

func (u authenticatedUser) GetUser() *database.User {
	return &u.user
}

// anonymousUser is used to represent a user that is not logged in.
type anonymousUser struct{}

func (u anonymousUser) IsAuthenticated() bool {
	return false
}

func (u anonymousUser) GetUser() *database.User {
	return nil
}

// GetUser returns the current user.
func GetUser(r *http.Request) User {
	// Get the session
	session := getUserSession(r)
	if session == nil {
		log.Debug("GetUser: no user session")
		return &anonymousUser{}
	}

	// Update LastSeen
	session.LastSeen = time.Now()
	database.DB.MustExec(
		"UPDATE user_session SET last_seen=$1 WHERE id=$2",
		session.LastSeen, session.ID)

	return &authenticatedUser{session.User}
}

// LogoutUser logs out the current user. It is safe to call this function if
// a user is not logged in.
func LogoutUser(w http.ResponseWriter, r *http.Request) {
	session := getUserSession(r)
	if session != nil {
		database.DB.MustExec("DELETE FROM user_session WHERE id=$1", session.ID)
	}
	cookie := &http.Cookie{Name: sessionKeyCookieName, MaxAge: -1}
	http.SetCookie(w, cookie)
}

// LoginUser logs in a user. InvalidUsernameOrPasswordError is returned if the
// username or password is incorrect.
func LoginUser(username, password string, w http.ResponseWriter) error {
	user := getUser(username, password)
	if user == nil {
		return InvalidUsernameOrPasswordError
	}
	sessionKey, err := createUserSession(*user)
	if err != nil {
		return err
	}
	cookie := &http.Cookie{Name: sessionKeyCookieName, Value: sessionKey}
	http.SetCookie(w, cookie)
	return nil
}

type userSession struct {
	database.User
	SessionID  uint
	SessionKey string
	LastSeen   time.Time
}

func getUserSession(r *http.Request) *userSession {
	sessionCookie, err := r.Cookie(sessionKeyCookieName)
	if err != nil {
		return nil
	}
	session := &userSession{}
	if err := database.DB.Get(session,
		`SELECT u.*, us.id AS session_id, us.key AS session_key, us.last_seen AS last_seen
		FROM user_session us
		JOIN user u ON us.user_id=u.id
		WHERE us.key=$1
		LIMIT 1`,
		sessionCookie.Value); err != nil {
		return nil
	}
	return session
}

func getUser(username, password string) *database.User {
	user := &database.User{}
	if err := database.DB.
		Get(user, "SELECT * FROM user WHERE username=$1 LIMIT 1", username); err != nil {
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
			"INSERT INTO user (username, password) VALUES ($1, $2)",
			username, passwordHash); err != nil {
			return err
		}
	}
	return nil
}

func usernameTaken(username string) bool {
	var user database.User
	err := database.DB.Get(&user, "SELECT * FROM user WHERE username=$1", username)
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
