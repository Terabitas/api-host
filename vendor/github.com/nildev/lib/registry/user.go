package registry

import (
	"encoding/json"
	"sync"
)

type (
	// User type
	User struct {
		data map[string]interface{}
		sync.RWMutex
	}

	// Tuple type
	Tuple struct {
		Key string
		Val interface{}
	}
)

// MakeUser constructor
func MakeUser(data map[string]interface{}) *User {
	return &User{
		data: data,
	}
}

// GetEmail returns use email
func (u *User) GetEmail() string {
	r, ok := u.Get("email")
	if ok {
		return r.(string)
	}

	return ""
}

// Set the given value under the specified key.
func (u *User) Set(key string, value interface{}) {
	// Get map shard.
	u.Lock()
	defer u.Unlock()
	u.data[key] = value
}

// Get an element from map under given key
func (u *User) Get(key string) (interface{}, bool) {
	// Get shard
	u.RLock()
	defer u.RUnlock()

	val, ok := u.data[key]
	return val, ok
}

// Remove an element from the map.
func (u *User) Remove(key string) {
	u.Lock()
	defer u.Unlock()
	delete(u.data, key)
}

// Count returns the number of elements within the map.
func (u *User) Count() int {
	return len(u.data)
}

// Has looks up an item under specified key
func (u *User) Has(key string) bool {
	u.RLock()
	defer u.RUnlock()

	_, ok := u.data[key]
	return ok
}

// Iterator returns an iterator which could be used in a for range loop.
func (u *User) Iterator() <-chan Tuple {
	ch := make(chan Tuple)
	go func() {
		u.RLock()
		for key, val := range u.data {
			ch <- Tuple{key, val}
		}
		u.RUnlock()
		close(ch)
	}()
	return ch
}

// MarshalJSON marshal "private" variables to json
func (u *User) MarshalJSON() ([]byte, error) {
	tmp := make(map[string]interface{})

	for item := range u.Iterator() {
		tmp[item.Key] = item.Val
	}
	return json.Marshal(tmp)
}
