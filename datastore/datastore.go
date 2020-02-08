package datastore

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Database struct {
	Name string
	Path string
	lock *FileLock
}

func NewDatabase(name string) *Database {

	// TODO: Should be Config ENV for persistent data
	path := "data"

	entry := &Database{
		Name: name,
		Path: path,
	}

	// Helpful mode
	entry.touch()

	return entry
}

var dir_perms = os.ModePerm

func (db *Database) touch() {
	os.MkdirAll(db.Path, dir_perms)
}

func (db *Database) filepath() string {
	return filepath.Join(db.Path, db.Name+".json")
}

// Locking is external to support transactional integrity
func (db *Database) Lock() {
	db.lock = Lock(db.filepath())
}

func (db *Database) Unlock() {
	db.lock.Unlock()
}

// TODO: Probably configurable as well ...
var open_perms os.FileMode = 0700

func (db *Database) Load(data interface{}) interface{} {
	fmt.Println("Loading Database")

	// Convert to JSON
	buffer := ReadFile(db.filepath())

	return Deserialize(buffer, data)
}

func (db *Database) Store(data interface{}) {
	fmt.Println("Storing Database")
	// Convert to Structure
	flattened := Serialize(data)
	WriteFile(db.filepath(), flattened)
}

func WriteFile(filepath string, data []byte) {
	err := ioutil.WriteFile(filepath, data, open_perms)
	if err != nil {
		fmt.Println("Can't read from file")
		panic("goodbye")
	}
}

func ReadFile(filepath string) []byte {
	buffer, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Println("Can't read from file")
		panic("goodbye")
	}

	return buffer
}
