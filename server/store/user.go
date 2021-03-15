package store

import (
	"fmt"

	"github.com/pkg/errors"
)

// UserStoreKeyPrefix prefix for user data key is KVStore.
const UserStoreKeyPrefix = "user_"

// UserStore API for user KVStore.
type UserStore interface {
	CurrentUser(mattermostUserID string) (SplunkUser, error)
	User(mattermostUserID string, server string, username string) (SplunkUser, error)

	ChangeCurrentUser(mattermostUserID string, userName string) error
	RegisterUser(mattermostUserID string, user SplunkUser) error
}

// SplunkUser stores splunk user info.
type SplunkUser struct {
	Server   string
	UserName string
	Token    string
}

// user KVStore value for each user
type user struct {
	LastLoginUserName string
	SplunkUsers       []SplunkUser
}

// CurrentUser returns last authorized user.
func (s *pluginStore) CurrentUser(mattermostUserID string) (SplunkUser, error) {
	su, err := s.loadUser(mattermostUserID)
	if err != nil {
		return SplunkUser{}, err
	}

	for _, u := range su.SplunkUsers {
		if u.UserName == su.LastLoginUserName {
			return u, nil
		}
	}
	return SplunkUser{}, errors.New("no user found")
}

// UserToken returns user info for given server and username
func (s *pluginStore) User(mattermostUserID string, server string, username string) (SplunkUser, error) {
	su, err := s.loadUser(mattermostUserID)
	if err != nil {
		return SplunkUser{}, err
	}

	for _, u := range su.SplunkUsers {
		if u.Server == server && u.UserName == username {
			return u, nil
		}
	}
	return SplunkUser{}, errors.New("no user found")
}

// ChangeCurrentUser changes authorized user to given one
// if userName is empty string it's equivalent of logout.
func (s *pluginStore) ChangeCurrentUser(mattermostUserID string, userName string) error {
	su, err := s.loadUser(mattermostUserID)
	if err != nil {
		return err
	}

	if userName == "" {
		goto found
	}

	for _, u := range su.SplunkUsers {
		if u.UserName == userName {
			goto found
		}
	}
	return errors.New("no user found")

found:
	su.LastLoginUserName = userName
	return s.storeUser(mattermostUserID, su)
}

// RegisterUser registers new splunk user
// if user with given username already exists than it will be overwritten.
func (s *pluginStore) RegisterUser(mattermostUserID string, user SplunkUser) error {
	su, err := s.loadUser(mattermostUserID)
	if err != nil {
		return err
	}

	ind := -1
	for i, u := range su.SplunkUsers {
		if u.UserName == user.UserName {
			ind = i
			break
		}
	}
	if ind != -1 {
		su.SplunkUsers = append(su.SplunkUsers[:ind], su.SplunkUsers[ind+1:]...)
	}

	su.SplunkUsers = append(su.SplunkUsers, user)
	return nil
}

func (s *pluginStore) loadUser(mattermostUserID string) (*user, error) {
	u := &user{}
	err := LoadGOB(s.userStore, fmt.Sprintf("%s%s", UserStoreKeyPrefix, mattermostUserID), u)
	if err != nil {
		return nil, errors.Wrapf(err, "error while loading a user with id : %s", mattermostUserID)
	}
	return u, nil
}

func (s *pluginStore) storeUser(mattermostUserID string, u *user) error {
	if u == nil {
		return errors.New("user is nil")
	}
	err := SetGOB(s.userStore, fmt.Sprintf("%s%s", UserStoreKeyPrefix, mattermostUserID), u)
	if err != nil {
		return errors.Wrap(err, "error while storing user")
	}
	return nil
}
