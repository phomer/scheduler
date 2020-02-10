package accounts

import (
	"errors"
	"log"
	"path/filepath"

	"github.com/phomer/scheduler/datastore"
)

type ClientConfig struct {
	Hostname    string
	Port        string
	Username    string
	Credentials string
}

var paths = []string{".", "~/.schedule"}

// Register initialization
func NewClientConfig(username string, hostname string) *ClientConfig {

	return &ClientConfig{
		Hostname: hostname,
		Username: username,
		Port:     "8000",
	}
}

// Client initialization
func FindClientConfig() *ClientConfig {
	file := FindFile(paths)
	buffer := datastore.ReadFile(file)

	config := new(ClientConfig)

	if len(buffer) > 0 {
		config = datastore.Deserialize(buffer, config).(*ClientConfig)
	} else {
		log.Fatal("Missing Config file", errors.New("Missing"), file)
	}

	return config
}

func (config *ClientConfig) filename() string {
	name := config.Hostname + "-" + config.Username + ".key"
	return name
}

func (config *ClientConfig) SaveConfig() {
	if config == nil {
		log.Fatal("Missing Config Data", nil)
	}

	buffer := datastore.Serialize(config)

	filename := config.filename()
	datastore.WriteFile(filepath.Join(".", filename), buffer)
}

func FindFile(paths []string) string {
	for i := 0; i < len(paths); i++ {
		path := "./"
		filename := "elephant-paulwhomer.key"
		datastore.TouchFile(path, filename)

		return filepath.Join(path, filename)
	}
	return ""
}
