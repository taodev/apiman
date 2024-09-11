package storage

import (
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

type KeyValueStore struct {
	data     map[string]any
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
	db.data = make(map[string]any)
	if err = yaml.Unmarshal(fileContent, &db.data); err != nil {
		db = nil
		return
	}

	return
}

func NewFromMemory() (db *KeyValueStore) {
	db = new(KeyValueStore)
	db.data = make(map[string]any)
	return
}

func (db *KeyValueStore) SetData(data map[string]any) {
	db.locker.Lock()
	defer db.locker.Unlock()

	db.data = data
}

func (db *KeyValueStore) Get(k string) (v any) {
	db.locker.RLock()
	defer db.locker.RUnlock()

	v = db.data[k]
	return
}

func (db *KeyValueStore) Set(k string, v any) {
	db.locker.Lock()
	defer db.locker.Unlock()

	db.data[k] = v
}

func (db *KeyValueStore) Del(k string) {
	db.locker.Lock()
	defer db.locker.Unlock()

	delete(db.data, k)
}

func (db *KeyValueStore) Each(f func(k string, v any)) {
	db.locker.RLock()
	defer db.locker.RUnlock()

	for k, v := range db.data {
		f(k, v)
	}
}

var EnvironmentDB *KeyValueStore
var GlobalDB *KeyValueStore = NewFromMemory()

func LoadEnv(name string) (err error) {
	// 判断env配置文件是否存在
	if _, err = os.Stat(name); err == nil {
		if EnvironmentDB, err = Load(name); err != nil {
			return
		}
	}

	return
}
