package datastore

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/phomer/scheduler/log"
)

// TODO: Should come from configuration
var open_perms os.FileMode = 0700
var dir_perms = os.ModePerm
var database_path = "data"

type Database struct {
	Name string
	Path string
	lock *FileLock
}

func NewDatabase(name string) *Database {

	entry := &Database{
		Name: name,
		Path: database_path,
	}

	// Make sure everything is initialized
	entry.Touch()

	return entry
}

// File lock on the underlying database data
func (db *Database) Lock() {
	db.lock = Lock(db.filepath())
}

// Release the file lock
func (db *Database) Unlock() {
	db.lock.Unlock()
}

func (db *Database) Load(data interface{}) interface{} {
	fmt.Println("Loading Database")

	// Convert to JSON
	buffer := ReadFile(db.filepath())

	return Deserialize(buffer, data)
}

func (db *Database) Store(data interface{}) {
	// Convert to Structure
	flattened := Serialize(data)

	WriteFile(db.filepath(), flattened)
}

func WriteFile(filepath string, data []byte) {
	err := ioutil.WriteFile(filepath, data, open_perms)
	if err != nil {
		log.Fatal("File", err, filepath)
	}
}

func ReadFile(filepath string) []byte {
	buffer, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Fatal("File", err, filepath)
	}

	return buffer
}

// See if the file exists, but throw an error when there are perm problems.
func FileExists(filepath string) bool {
	// Dela with the file
	_, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		return false

	} else if err != nil {
		log.Fatal("Reaching File", err)
	}

	return true
}

func TouchFile(path string, filename string) {
	os.MkdirAll(path, dir_perms)

	file := filepath.Join(path, filename)

	if !FileExists(file) {
		empty, err := os.Create(file)
		if err != nil {
			log.Fatal("Creating File", err)
		}
		empty.Close()
	}
}

func (db *Database) Touch() {
	TouchFile(db.Path, db.filename())
}

func (db *Database) filename() string {
	return db.Name + ".json"
}

func (db *Database) filepath() string {
	return filepath.Join(db.Path, db.filename())
}
