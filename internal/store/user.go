package store

import (
	"github.com/sirupsen/logrus"
	"sync"
)

type user struct {
	user  map[string]string
	mutex sync.RWMutex
}

func NewUser() *user {
	return &user{
		user: map[string]string{},
	}
}

func (u *user) Add(username, password string) {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	logrus.WithField("user", username).Info("Add user.")
	u.user[username] = password
}

func (u *user) IsValid(username, password string) bool {
	u.mutex.RLock()
	defer u.mutex.RUnlock()

	if pw, ok := u.user[username]; ok {
		return pw == password
	}

	return true
}
