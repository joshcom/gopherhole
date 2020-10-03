package main

import (
	"flag"
	"github.com/joshcom/gopherhole/gopherhole"
	"log"
)

var (
	help       = flag.Bool("help", false, "Prints this help message.")
	configfile = flag.String("config", "", "The path the configuration file.")
	root       = flag.String("root", "", "The path to the root directory to serve files from.")
	host       = flag.String("host", "", "The host name clients will use when making requests to this gopher hole.")
	port       = flag.Int("port", 0, "The port to listen on.")
	maxconn    = flag.Int("maxconn", -1, "The maximum number of concurrent connections. 0 means no maximum.")
	mapfile    = flag.String("mapfile", "", "Files with this name will be served instead of the contents of the directory they reside in.")
	nologging  = flag.Bool("nologging", false, "Logging output will be disabled.")
	logfile    = flag.String("logfile", "", "Logs will be printed here instead of the STDOUT.")
	timeout    = flag.Int("timeout", -1, "After waiting this many seconds for a query, connections will be closed. 0 means no timeout.")
)

func main() {
	flag.Parse()

	if *help {
		flag.PrintDefaults()
		return
	}

	config := buildConfig()
	server := gopherhole.NewServer(config)
	server.Run()
}

func buildConfig() (config gopherhole.Configuration) {
	var err error

	if *configfile != "" {
		config, err = gopherhole.NewConfigurationFromFile(*configfile)
		if err != nil {
			log.Fatalf("Failed to open configuration file '%s.", *configfile)
		}
	} else {
		config = gopherhole.NewConfiguration()
	}

	if *root != "" {
		config.RootDirectory = *root
	}

	if *host != "" {
		config.Host = *host
	}

	if *port != 0 {
		config.Port = *port
	}

	if *maxconn >= 0 {
		config.MaxConnections = *maxconn
	}

	if *mapfile != "" {
		config.MapFileName = *mapfile
	}

	if *nologging {
		config.LogDisabled = *nologging
	}

	if *logfile != "" {
		config.LogFile = *logfile
	}

	if *timeout >= 0 {
		config.IdleTimeout = *timeout
	}

	return
}
