package store

import (
	"encoding/json"
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
	DeleteUser(mattermostUserID string, server string, userName string) error
	GetSubscription(key string) (map[string][]string, error)
	SetSubscription(key string, subscription map[string][]string) error
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
func (s *pluginStore) RegisterUser(mattermostUserID string, splunkUser SplunkUser) error {
	su, err := s.loadUser(mattermostUserID)
	if err != nil {
		su = &user{
			LastLoginUserName: "",
			SplunkUsers:       []SplunkUser{},
		}
	}

	ind := -1
	for i, u := range su.SplunkUsers {
		if u.UserName == splunkUser.UserName {
			ind = i
			break
		}
	}
	if ind != -1 {
		su.SplunkUsers = append(su.SplunkUsers[:ind], su.SplunkUsers[ind+1:]...)
	}

	su.SplunkUsers = append(su.SplunkUsers, splunkUser)
	return s.storeUser(mattermostUserID, su)
}

// DeleteUser deletes user from KV store
func (s *pluginStore) DeleteUser(mattermostUserID string, server string, userName string) error {
	su, err := s.loadUser(mattermostUserID)
	if err != nil {
		return errors.Wrap(err, "no user found")
	}

	ind := -1
	for i, u := range su.SplunkUsers {
		if u.Server == server && u.UserName == userName {
			ind = i
			break
		}
	}
	if ind != -1 {
		su.SplunkUsers = append(su.SplunkUsers[:ind], su.SplunkUsers[ind+1:]...)
	}

	return s.storeUser(mattermostUserID, su)
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

func (s *pluginStore) GetSubscription(key string) (map[string][]string, error) {
	var subscriptions map[string][]string
	subscriptionByte, appErr := s.userStore.Load(key)
	if appErr != nil {
		return subscriptions, errors.Wrap(appErr, "Error While Getting Subscription From KV Store")
	}
	if subscriptionByte != nil {
		err := json.Unmarshal(subscriptionByte, &subscriptions)
		if err != nil {
			return subscriptions, errors.Wrap(err, "Error While Decoding The Subscriptions")
		}
	}
	return subscriptions, nil
}

func (s *pluginStore) SetSubscription(key string, subscription map[string][]string) error {
	value, err := json.Marshal(subscription)
	if err != nil {
		return errors.Wrap(err, "Error While Encoding The Subscriptions")
	}
	return s.userStore.Store(key, value)
}
