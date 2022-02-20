package engine

import (
	"errors"
	"sync"
)

type ContentType = string
type CharSet = string
type NabiaRecord struct {
	RawData     []byte
	ContentType ContentType // "Content-Type" https://datatracker.ietf.org/doc/html/rfc2616/#section-14.17
	//RWMutex     sync.RWMutex
}

func NewNabiaString(s string) *NabiaRecord {
	return &NabiaRecord{RawData: []byte(s), ContentType: "text/plain; charset=UTF-8"}
}

func NewNabiaRecord(data []byte, ct ContentType) *NabiaRecord {
	return &NabiaRecord{RawData: data, ContentType: ct}
}

type path = string
type NabiaDB struct {
	Records map[path]*NabiaRecord // Key = path; value = pointer to content
	RWMutex sync.RWMutex
}

func NewNabiaDB() *NabiaDB {

	return &NabiaDB{Records: map[string]*NabiaRecord{}, RWMutex: sync.RWMutex{}}
}

// Below are the DB primitives.

// Exists checks if the key name provided exists in the Nabia map. It locks
// to read and unlocks immediately after.
func (ns *NabiaDB) Exists(key string) bool {
	var exists bool
	exists = false
	ns.RWMutex.RLock()
	defer ns.RWMutex.RUnlock()
	if ns.Records != nil && ns.Records[key] != nil {
		exists = true
	}
	return exists
}

// Read takes a key name and attempts to pull the data from the Nabia DB map.
// Returns a NabiaRecord (if found) and an error (if not found). Callers must
// always check the error returned in the second parameter, as the result cannot
// be used if the "error" field is not nil. This function is safe to call even
// with empty data, because the method applies a mutex, checks for the existence
// of the record, and if and only if this record exists, it returns it before
// unlocking. Otherwise, it returns an error and unlocks.
func (ns *NabiaDB) Read(key string) (NabiaRecord, error) {
	var nr NabiaRecord
	var err error
	ns.RWMutex.RLock()
	defer ns.RWMutex.RUnlock()
	if ns.Exists(key) {
		nr = *ns.Records[key]
	} else {
		err = errors.New("can't read a key that doesn't exist")
	}
	return nr, err
}

// Write takes the key and a value of NabiaRecord datatype and places it on the
// database, potentially overwriting whatever was there before, because Write
// has no data safety features preventing the overwriting of data.
func (ns *NabiaDB) Write(key string, value NabiaRecord) {
	ns.RWMutex.Lock()
	defer ns.RWMutex.Unlock()
	ns.Records[key] = &value
}

// Destroy takes a key and removes it from the map. This method doesn't have
// existence-checking logic. It is safe to use on empty data, it simply doesn't
// do anything if the record doesn't exist. Internally, an empty record looks
// like: (record[key] = NabiaData{}) where the NabiaData is a pointer to nil.
func (ns *NabiaDB) Destroy(key string) { // Doesn't check if exists
	ns.RWMutex.Lock()
	defer ns.RWMutex.Unlock()
	delete(ns.Records, key)
}
