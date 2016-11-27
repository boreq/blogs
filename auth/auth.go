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
	user *database.User
}

func (u *authenticatedUser) IsAuthenticated() bool {
	return true
}

func (u *authenticatedUser) GetUser() *database.User {
	return u.user
}

// anonymousUser is used to represent a user that is not logged in.
type anonymousUser struct{}

func (u *anonymousUser) IsAuthenticated() bool {
	return false
}

func (u *anonymousUser) GetUser() *database.User {
	return nil
}

// GetUser returns the current user.
func GetUser(r *http.Request) User {
	// Get the session
	session := getUserSession(r)
	if session == nil {
		log.Debug("Session doesn't exist")
		return &anonymousUser{}
	}
	// Update LastSeen
	session.LastSeen = time.Now()
	database.DB.Save(session)
	// Get the user
	user := &database.User{}
	if err := database.DB.
		Model(&session).
		Related(&user).
		Error; err != nil {
		log.Debug("User doesn't exist")
		return &anonymousUser{}
	}
	return &authenticatedUser{user}
}

// LogoutUser logs out the current user. It is safe to call this function if
// a user is not logged in.
func LogoutUser(w http.ResponseWriter, r *http.Request) {
	session := getUserSession(r)
	if session != nil {
		database.DB.Delete(&session)
	}
	cookie := &http.Cookie{Name: sessionKeyCookieName, MaxAge: -1}
	http.SetCookie(w, cookie)
}

// LoginUser logs in a specified user. Returns InvalidUsernameOrPasswordError
// if the usrename and password pair is incorrect or other errors.
func LoginUser(username, password string, w http.ResponseWriter) error {
	user := getUser(username, password)
	if user == nil {
		return InvalidUsernameOrPasswordError
	}
	sessionKey, err := createUserSession(user)
	if err != nil {
		return err
	}
	cookie := &http.Cookie{Name: sessionKeyCookieName, Value: sessionKey}
	http.SetCookie(w, cookie)
	return nil
}

func getUserSession(r *http.Request) *database.UserSession {
	sessionCookie, err := r.Cookie(sessionKeyCookieName)
	if err != nil {
		return nil
	}
	session := &database.UserSession{}
	if database.DB.
		Where(database.UserSession{SessionKey: sessionCookie.Value}).
		First(session).
		RecordNotFound() {
		return nil
	}
	return session
}

func getUser(username, password string) *database.User {
	user := &database.User{}
	if database.DB.
		Where(database.User{Username: username}).
		First(user).
		RecordNotFound() {
		return nil
	}
	if compareHashAndPassword(user.Password, password) {
		return user
	}
	return nil
}

func createUserSession(user *database.User) (string, error) {
	sessionKey := make([]byte, sessionKeySize)
	if err := generateSessionKey(sessionKey); err != nil {
		return "", err
	}
	sessionKeyString := fmt.Sprintf("%x", sessionKey)
	session := database.UserSession{
		UserID:     user.ID,
		SessionKey: sessionKeyString,
		LastSeen:   time.Now(),
	}
	if err := database.DB.Create(&session).Error; err != nil {
		return "", err
	}

	return sessionKeyString, nil

}

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
		user := database.User{
			Username: username,
			Password: passwordHash,
		}
		if err := database.DB.Create(&user).Error; err != nil {
			return err
		}
	}
	return nil
}

func usernameTaken(username string) bool {
	user := &database.User{}
	return !database.DB.
		Where(database.User{Username: username}).
		First(&user).
		RecordNotFound()
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
