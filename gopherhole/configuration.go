package gopherhole

import (
	"encoding/json"
	"os"
)

const (
	ROOT_DIRECTORY    = "/var/gopherhole"
	PORT              = 70
	MAX_CONNECTIONS   = 0
	HOST              = "localhost"
	MAP_FILE_NAME     = "gophermap"
	LOG_DISABLED      = false
	LOG_FILE          = ""
	IDLE_TIMEOUT      = 60
	DEFAULT_MIME_TYPE = "text/plain"
)

type Configuration struct {
	RootDirectory      string
	Host               string
	Port               int
	MaxConnections     int
	MapFileName        string
	LogDisabled        bool
	LogFile            string
	IdleTimeout        int
	DefaultMimeType    string
	MimeTypeIgnoreList []string
}

func NewConfiguration() Configuration {
	return Configuration{
		RootDirectory:   ROOT_DIRECTORY,
		Host:            HOST,
		Port:            PORT,
		MaxConnections:  MAX_CONNECTIONS,
		MapFileName:     MAP_FILE_NAME,
		LogDisabled:     LOG_DISABLED,
		LogFile:         LOG_FILE,
		IdleTimeout:     IDLE_TIMEOUT,
		DefaultMimeType: DEFAULT_MIME_TYPE,
	}
}

func NewConfigurationFromFile(pathToFile string) (config Configuration, err error) {
	file, err := os.Open(pathToFile)
	if err != nil {
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	config = NewConfiguration()
	err = decoder.Decode(&config)

	return
}
