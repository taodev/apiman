package storage

import (
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

type KeyValueStore struct {
	data     map[string]string
	locker   sync.RWMutex
	filepath string
}

func Load(filepath string) (db *KeyValueStore, err error) {
	fileContent, err := os.ReadFile(filepath)
	if err != nil {
		return
	}

	db = new(KeyValueStore)
	db.filepath = filepath
	db.data = make(map[string]string)
	if err = yaml.Unmarshal(fileContent, &db.data); err != nil {
		db = nil
		return
	}

	return
}

func NewFromMemory() (db *KeyValueStore) {
	db = new(KeyValueStore)
	db.data = make(map[string]string)
	return
}

func (db *KeyValueStore) Get(k string) (v string) {
	db.locker.RLock()
	defer db.locker.RUnlock()

	v = db.data[k]
	return
}

func (db *KeyValueStore) Set(k string, v string) {
	db.locker.Lock()
	defer db.locker.Unlock()

	db.data[k] = v
}

func (db *KeyValueStore) Del(k string) {
	db.locker.Lock()
	defer db.locker.Unlock()

	delete(db.data, k)
}

func (db *KeyValueStore) Each(f func(k, v string)) {
	db.locker.RLock()
	defer db.locker.RUnlock()

	for k, v := range db.data {
		f(k, v)
	}
}

var EnvironmentDB *KeyValueStore
var GlobalDB *KeyValueStore = NewFromMemory()
