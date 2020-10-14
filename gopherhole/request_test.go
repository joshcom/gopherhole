package gopherhole

import (
	"strings"
	"testing"
)

func TestRequest_newRequest(t *testing.T) {
	config := buildConfiguration()
	req := newRequest("phlog/", config)

	if req.selector.query != "phlog/" ||
		req.configuration.RootDirectory != config.RootDirectory {
		t.Error("Request object not set up correctly.")
	}
}

func TestRequest_process(t *testing.T) {
	config := buildConfiguration()
	config.MimeTypeIgnoreList = []string{"application/"}

	t.Run("process file", func(t *testing.T) {
		buffer, writer := newBytesWriter()
		req := newRequest("art/laptop.txt", config)
		_, err := req.process(writer)
		if err != nil {
			t.Fatalf("Unexpected error %v", err)
		}

		dataStr := buffer.String()
		if strings.Index(dataStr, "$ vim") < 0 {
			t.Error("Payload data not as expected.")
		}
	})

	t.Run("process directory with mapfile", func(t *testing.T) {
		buffer, writer := newBytesWriter()
		req := newRequest("phlog/", config)
		_, err := req.process(writer)
		if err != nil {
			t.Fatalf("Unexpected error %v", err)
		}

		dataStr := buffer.String()
		expected := "ijoshcom.net - PHLOG\t(NOTHING)\tnohost\t0\r\n"
		if strings.Index(dataStr, expected) < 0 {
			t.Error("Payload data not as expected.")
		}
	})

	t.Run("process directory without mapfile", func(t *testing.T) {
		buffer, writer := newBytesWriter()
		req := newRequest("files", config)
		_, err := req.process(writer)
		if err != nil {
			t.Fatalf("Unexpected error %v", err)
		}

		dataStr := buffer.String()
		if strings.Index(dataStr, "happydance.gif") < 0 {
			t.Error("Payload data not as expected.")
		}

		if strings.Index(dataStr, "run.exe") >= 0 {
			t.Error("Payload data not as expected.")
		}
	})

	t.Run("selector error", func(t *testing.T) {
		buffer, writer := newBytesWriter()
		req := newRequest("../../", config)
		_, err := req.process(writer)
		if err != nil {
			t.Error("Disallowed path should be handled without error.")
		}

		dataStr := buffer.String()
		if strings.Index(dataStr, "3Invalid query") < 0 {
			t.Error("Error payload not as expected.")
		}
	})

	t.Run("resource error", func(t *testing.T) {
		buffer, writer := newBytesWriter()
		req := newRequest("iheartsocialmedia.txt", config)
		_, err := req.process(writer)
		if err != nil {
			t.Error("Invalid file should be handled without error.")
		}
		dataStr := buffer.String()
		if strings.Index(dataStr, "3File not found.") < 0 {
			t.Error("Error payload not as expected.")
		}
	})

	t.Run("restricted resource error", func(t *testing.T) {
		buffer, writer := newBytesWriter()
		req := newRequest("files/run.exe", config)
		_, err := req.process(writer)
		if err != nil {
			t.Error("Invalid file should be handled without error.")
		}
		dataStr := buffer.String()
		if strings.Index(dataStr, "3File not found.") < 0 {
			t.Error("Error payload not as expected.")
		}
	})
}
