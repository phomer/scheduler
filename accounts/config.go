package accounts

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/phomer/scheduler/datastore"
	"github.com/phomer/scheduler/log"
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
	file, err := FindFile(paths)
	if err != nil {
		log.Fatal("Config File", err)
	}

	buffer := datastore.ReadFile(file)

	config := new(ClientConfig)

	if len(buffer) > 0 {
		config = datastore.Deserialize(buffer, config).(*ClientConfig)
	} else {
		log.Fatal("Missing Config file", errors.New("Missing"), file)
	}

	return config
}

func (config *ClientConfig) GetUrl(request_type string) string {

	hostname := config.Hostname

	if Hostname() == hostname {
		hostname = "127.0.0.1"
	}

	fmt.Println("Accessing server:", hostname)

	return fmt.Sprintf("%s://%s:%s/%s", config.Protocol, hostname, config.Port, request_type)
}

func (config *ClientConfig) filename() string {
	return fmt.Sprintf("%s-%s.key", config.Hostname, config.Username)
}

func (config *ClientConfig) SaveConfig() {
	if config == nil {
		log.Fatal("Missing Config Data", nil)
	}

	buffer := datastore.Serialize(config)

	filename := config.filename()
	datastore.WriteFile(filepath.Join(".", filename), buffer)
}

func FindFile(paths []string) (string, error) {

	for i := 0; i < len(paths); i++ {
		matches, err := filepath.Glob("*-*.key")
		if err != nil {
			return "", err
		}

		// TODO: Should realy ask the user which one they want, but for
		// now lets just take the first file that is found.
		if len(matches) > 0 {
			return matches[0], nil
		}
	}
	return "", errors.New("Missing File")
}
