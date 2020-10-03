package gopherhole

import (
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
	"testing"
)

func TestServer_Run(t *testing.T) {
	listener, resetFunc := buildTestListener()
	defer resetFunc()

	server := buildServer()
	go func() { server.Run() }()
	defer listener.Close()

	t.Run("test successful query", func(t *testing.T) {
		listener.clientConn.Write([]byte("phlog/\r\n"))

		var readData []byte
		readData, err := ioutil.ReadAll(listener.clientConn)
		if err != nil {
			t.Fatalf("Error reading response from server: %v", err)
		}

		expected := "1^^Top of the gopherhole^^\t/"
		if strings.Index(string(readData), expected) < 0 {
			t.Error("Payload not returned as expected.")
		}
	})
}

func TestServer_Terminate(t *testing.T) {
	listener, resetFunc := buildTestListener()
	defer resetFunc()

	server := buildServer()
	go func() { server.Run() }()
	defer listener.Close()

	server.Terminate()
	num, err := listener.clientConn.Write([]byte("phlog/\r\n"))

	if err == nil || num != 0 {
		t.Errorf("Expected failed write to server after termination.")
	}
}

func TestServer_HandleConnection(t *testing.T) {
	listener := newStubListener()
	logger := log.New(os.Stdout, "[test]", 0)
	logger.SetOutput(ioutil.Discard)
	server := buildServer()
	defer listener.Close()

	var err error
	var readData []byte

	t.Run("test successful query", func(t *testing.T) {
		go func() {
			server.handleConnection(&listener.serverConn, logger)
		}()

		_, err = listener.clientConn.Write([]byte("phlog/\r\n"))
		readData, err = ioutil.ReadAll(listener.clientConn)

		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		expected := "1^^Top of the gopherhole^^\t/"
		if strings.Index(string(readData), expected) < 0 {
			t.Error("Payload not returned as expected.")
		}
	})
}

type stubListener struct {
	serverConn net.Conn
	clientConn net.Conn
}

func newStubListener() stubListener {
	clientP, serverP := net.Pipe()
	return stubListener{
		clientConn: clientP,
		serverConn: serverP,
	}
}

func (l *stubListener) Accept() (net.Conn, error) {
	return l.serverConn, nil
}

func (l *stubListener) Close() (err error) {
	l.serverConn.Close()
	l.clientConn.Close()
	return
}

func (l *stubListener) Addr() (addr net.Addr) {
	return
}

func buildTestListener() (stubListener, func()) {
	listener := newStubListener()
	originalListen := Listen

	resetFunc := func() {
		Listen = originalListen
		listener.Close()
	}

	Listen = func(network, address string) (l net.Listener, err error) {
		l = &listener
		return l, nil
	}

	return listener, resetFunc
}

func buildServer() (server Server) {
	config := NewConfiguration()
	config.RootDirectory = "testdata/mygopherhole"
	config.IdleTimeout = 10
	config.LogDisabled = true

	server = NewServer(config)
	return server
}
