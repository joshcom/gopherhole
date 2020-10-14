/*
Package gopherhole implements a gopher server.


Example usage:

  // Any omitted config setting will use
  // the default noted below.
  // Prefer absolute paths.
  // config := gopherhole.NewConfiguration() loads the defaults.
  //
  // View the README for more details on configuration options.
  config := gopherhole.Configuration{
    RootDirectory:  "/var/gopherhole",
    Host:           "localhost",
    Port:           70,
    MaxConnections: 0, // 0 means no maximum.
    MapFileName:    "gophermap",
    LogFile:        "", // An empty string prints to STDOUT
    IdleTimeout:    60, // The number of seconds a client has after connecting
                        // to make a query.
  }
  server := gopherhole.NewServer(config)
  server.Run()

Additionally, the configuration can be stored as JSON
in a file:

  {
    "Host": "localhost",
    // ...
  }

The absolute path can be specified when building a new configuration:

  config := gopherhole.NewConfigurationFromFile("/path/to/config.json")

*/
package gopherhole

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"os"
	"sync"
	"time"
)

const (
	CONNECTION_FREED_BUFFER_SIZE = 10
)

type Server struct {
	Configuration Configuration
	Logger        *log.Logger
	handlerCount  int

	// If a message is sent to this channel, the server terminates.
	term chan bool

	// When the connection limit is reached, this channel will recieve
	// a message when a connection is freed.
	free chan bool

	// Connection counter mutex.
	mu *sync.Mutex
}

func NewServer(config Configuration) (server Server) {
	server = Server{Configuration: config}
	server.term = make(chan bool)
	server.free = make(chan bool, CONNECTION_FREED_BUFFER_SIZE)
	server.handlerCount = 0
	server.mu = new(sync.Mutex)

	return server
}

func (s *Server) Run() (err error) {
	f := os.Stdout
	if s.Configuration.LogDisabled {
		log.SetFlags(0)
		log.SetOutput(ioutil.Discard)
	} else if s.Configuration.LogFile != "" {
		f, err = os.OpenFile(s.Configuration.LogFile,
			os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("Unable to open logfile %s: %v", s.Configuration.LogFile, err)
		}
		defer f.Close()
	}

	port := fmt.Sprintf(":%d", s.Configuration.Port)
	ln, err := Listen("tcp", port)
	if err != nil {
		log.Fatalf("Error listening to %s, %s", port, err)
		return
	}
	defer ln.Close()

	for {
		if s.handlersAtCapacity() {
			select {
			case <-s.free:
				continue
			}
		}

		conn, err := ln.Accept()
		logger := s.buildLogger(f)
		if err != nil {
			logger.Printf("Accept failure: %v", err)
			continue
		} else {
			logger.Printf("Accepting connection from %s",
				conn.RemoteAddr().String())

			s.incrementHandlerCounter()
			go s.handleConnection(&conn, logger)
		}
	}
}

func (s *Server) handleConnection(c *net.Conn, logger *log.Logger) {
	defer (*c).Close()
	defer s.decrementHandlerCounter()

	var ch = make(chan error)
	var err error
	var clientData string

	go func(c *net.Conn) {
		clientData, err = bufio.NewReader(*c).ReadString('\n')
		ch <- err
	}(c)

	duration := time.Duration(s.Configuration.IdleTimeout)
	select {
	case err = <-ch:
		// _Something_ happened...
	case <-time.After(duration * time.Second):
		logger.Printf("Timed out waiting %d seconds",
			s.Configuration.IdleTimeout)
		return
	}

	if err != nil {
		logger.Printf("Failed to read data from client:\n%v", err)
		return
	}

	logger.Printf("Query: %s", clientData)
	request := newRequest(clientData, s.Configuration)

	writer := io.Writer(*c)
	num, err := request.process(&writer)
	if err != nil {
		logger.Printf("Failed to process request:\n%v", err)
		return
	}

	logger.Printf("Sent %d bytes.", num)
	return
}

func (s *Server) incrementHandlerCounter() {
	s.mu.Lock()
	s.handlerCount += 1
	s.mu.Unlock()

	if s.handlersAtCapacity() {
		log.Printf("Concurrent request capacity reached (%d).",
			s.Configuration.MaxConnections)
	}
}

func (s *Server) decrementHandlerCounter() {
	s.mu.Lock()
	if s.handlersAtCapacity() {
		s.handlerCount -= 1
		s.free <- true
	} else {
		s.handlerCount -= 1
	}
	s.mu.Unlock()
}

func (s *Server) handlersAtCapacity() bool {
	if s.Configuration.MaxConnections <= 0 {
		return false
	}

	return s.handlerCount >= s.Configuration.MaxConnections
}

func (s *Server) buildLogger(f *os.File) (logger *log.Logger) {
	loggerID, err := s.generateID()
	if err != nil {
		log.Printf("Failed to generate unique ID for logger: %v", err)
		log.Print("Using 'noid' for request.")
		loggerID = "noid"
	}

	logger = log.New(f, loggerID, log.LstdFlags)
	if s.Configuration.LogDisabled {
		logger.SetOutput(ioutil.Discard)
		logger.SetFlags(0)
	}

	return logger
}

func (s *Server) generateID() (id string, err error) {
	randBytes := make([]byte, 16)
	_, err = rand.Read(randBytes)
	if err != nil {
		log.Printf("Error generating logging id: %v", err)
		return
	}
	id = fmt.Sprintf("[%x-%x-%x-%x-%x] ",
		randBytes[0:4],
		randBytes[4:6],
		randBytes[6:8],
		randBytes[8:10],
		randBytes[10:])
	return
}

var Listen = func(network, address string) (net.Listener, error) {
	ln, err := net.Listen(network, address)
	return ln, err
}
