package gopherhole

import (
	"strings"
	"testing"
)

func buildConfiguration() (config Configuration) {
	config = NewConfiguration()
	config.RootDirectory = "testdata/mygopherhole"
	return config
}

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
		req := newRequest("art/laptop.txt", config)
		err := req.process()
		if err != nil {
			t.Fatalf("Unexpected error %v", err)
		}

		dataStr := string(*req.payload)
		if strings.Index(dataStr, "$ vim") < 0 {
			t.Error("Payload data not as expected.")
		}
	})

	t.Run("process directory with mapfile", func(t *testing.T) {
		req := newRequest("phlog/", config)
		err := req.process()
		if err != nil {
			t.Fatalf("Unexpected error %v", err)
		}

		dataStr := string(*req.payload)
		expected := "ijoshcom.net - PHLOG\t(NOTHING)\tnohost\t0\r\n"
		if strings.Index(dataStr, expected) < 0 {
			t.Error("Payload data not as expected.")
		}
	})

	t.Run("process directory without mapfile", func(t *testing.T) {
		req := newRequest("files", config)
		err := req.process()
		if err != nil {
			t.Fatalf("Unexpected error %v", err)
		}

		dataStr := string(*req.payload)
		if strings.Index(dataStr, "happydance.gif") < 0 {
			t.Error("Payload data not as expected.")
		}

		if strings.Index(dataStr, "run.exe") >= 0 {
			t.Error("Payload data not as expected.")
		}
	})

	t.Run("selector error", func(t *testing.T) {
		req := newRequest("../../", config)
		err := req.process()
		if err != nil {
			t.Error("Disallowed path should be handled without error.")
		}

		dataStr := string(*req.payload)
		if strings.Index(dataStr, "3Invalid query") < 0 {
			t.Error("Error payload not as expected.")
		}
	})

	t.Run("resource error", func(t *testing.T) {
		req := newRequest("iheartsocialmedia.txt", config)
		err := req.process()
		if err != nil {
			t.Error("Invalid file should be handled without error.")
		}
		dataStr := string(*req.payload)
		if strings.Index(dataStr, "3Requested resource not found.") < 0 {
			t.Error("Error payload not as expected.")
		}
	})
}
