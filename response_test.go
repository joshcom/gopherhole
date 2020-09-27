package gopherhole

import (
	"strings"
	"testing"
)

func TestResponse_build(t *testing.T) {
	resp := response{
		host:               "joshcom.net",
		port:               70,
		rootDirectory:      "testdata/mygopherhole",
		mapFileName:        "gophermap",
		mimeTypeIgnoreList: []string{"application/"},
	}

	t.Run("process file", func(t *testing.T) {
		path := "testdata/mygopherhole/art/laptop.txt"
		res, _ := newResource(path)
		data, err := resp.build(&res)

		if err != nil {
			t.Fatalf("Unexpected error %v", err)
		}

		dataStr := string(*data)
		if strings.Index(dataStr, "$ vim") < 0 {
			t.Error("Payload data not as expected.")
		}
	})

	t.Run("process directory with mapfile", func(t *testing.T) {
		path := "testdata/mygopherhole/phlog/"
		res, _ := newResource(path)
		data, err := resp.build(&res)

		if err != nil {
			t.Fatalf("Unexpected error %v", err)
		}

		dataStr := string(*data)
		expected := "ijoshcom.net - PHLOG\t(NOTHING)\tnohost\t0\r\n"
		if strings.Index(dataStr, expected) < 0 {
			t.Error("Payload data not as expected.")
		}
	})

	t.Run("process directory without mapfile", func(t *testing.T) {
		path := "testdata/mygopherhole/files"
		res, _ := newResource(path)
		data, err := resp.build(&res)

		if err != nil {
			t.Fatalf("Unexpected error %v", err)
		}

		dataStr := string(*data)
		if strings.Index(dataStr, "happydance.gif") < 0 {
			t.Error("Payload data not as expected.")
		}
		if strings.Index(dataStr, "run.exe") >= 0 {
			t.Error("Payload data not as expected.")
		}
	})

	t.Run("resource error", func(t *testing.T) {
		path := "testdata/mygopherhole/art/iheartsocialmedia.txt"
		res, _ := newResource(path)
		data, err := resp.build(&res)

		if err != nil {
			t.Error("Invalid file should be handled without error.")
		}
		dataStr := string(*data)
		if strings.Index(dataStr, "3Requested resource not found.") < 0 {
			t.Error("Error payload not as expected.")
		}
	})
}
