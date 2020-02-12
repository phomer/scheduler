package accounts

import (
	"errors"
	"fmt"
	"log"
	"path/filepath"

	"github.com/phomer/scheduler/datastore"
)

type ClientConfig struct {
	Protocol    string
	Hostname    string
	Port        string
	Username    string
	Credentials string
	Token       *Token
}

var paths = []string{".", "~/.schedule"}

// Register initialization
func NewClientConfig(hostname string, username string, token *Token) *ClientConfig {

	return &ClientConfig{
		Protocol: "http",
		Hostname: hostname,
		Username: username,
		Port:     "8000",
		Token:    token,
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

func (config *ClientConfig) GetUrl() string {

	hostname := config.Hostname

	// TODO: Replace with DNS lookup?
	if Hostname() == hostname {
		hostname = "127.0.0.1"
	}

	fmt.Println("Hostname ", hostname, " and ", Hostname())

	return fmt.Sprintf("%s://%s:%s/", config.Protocol, hostname, config.Port)
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
