package store

import (
	"bytes"
	"encoding/gob"
	"encoding/json"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/pkg/errors"
)

// API that store uses for interactions with KVStore
type API interface {
	KVGet(key string) ([]byte, *model.AppError)
	KVSet(key string, value []byte) *model.AppError
	KVDelete(key string) *model.AppError
	LogDebug(msg string, keyValuePairs ...interface{})
	LogInfo(msg string, keyValuePairs ...interface{})
	LogError(msg string, keyValuePairs ...interface{})
	LogWarn(msg string, keyValuePairs ...interface{})
}

// KVStore abstraction for plugin.API.KVStore
type KVStore interface {
	Load(key string) ([]byte, error)
	Store(key string, data []byte) error
	Delete(key string) error
	loadGOB(key string, v interface{}) error
	setGOB(key string, v interface{}) error
	setJSON(key string, v interface{}) error
	loadJSON(key string, v interface{}) error
}

type store struct {
	api API
}

// NewStore creates KVStore from plugin.API
func NewStore(api API) KVStore {
	return &store{
		api: api,
	}
}

func (s *store) Load(key string) ([]byte, error) {
	data, appErr := s.api.KVGet(key)
	if appErr != nil {
		return nil, errors.WithMessage(appErr, "failed plugin KVGet")
	}

	return data, nil
}

func (s *store) Store(key string, data []byte) error {
	appErr := s.api.KVSet(key, data)
	if appErr != nil {
		return errors.Wrapf(appErr, "Error while storing data with KVStore with key : %q", key)
	}
	return nil
}

func (s *store) Delete(key string) error {
	appErr := s.api.KVDelete(key)
	if appErr != nil {
		return errors.Wrapf(appErr, "Error while deleting data from KVStore with key : %q", key)
	}
	return nil
}

// loadGOB load json data from KVStore
func (s *store) loadGOB(key string, v interface{}) (returnErr error) {
	data, err := s.Load(key)
	if err != nil {
		return errors.Wrap(err, "Error while loading json")
	}
	return gob.NewDecoder(bytes.NewBuffer(data)).Decode(v)
}

// setGOB sets json data in KVStore
func (s *store) setGOB(key string, v interface{}) (returnErr error) {
	data := bytes.Buffer{}
	err := gob.NewEncoder(&data).Encode(v)
	if err != nil {
		return errors.Wrap(err, "Error while storing json")
	}
	return s.Store(key, data.Bytes())
}

func (s *store) loadJSON(key string, v interface{}) (returnErr error) {
	bytes, err := s.Load(key)
	if err != nil {
		return err
	}
	if bytes == nil {
		return nil
	}
	return json.Unmarshal(bytes, v)
}

func (s *store) setJSON(key string, v interface{}) error {
	bytes, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return s.Store(key, bytes)
}
