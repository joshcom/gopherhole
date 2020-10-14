package gopherhole

import (
	"strings"
	"testing"
)

func TestResponse_write(t *testing.T) {
	resp := response{
		host:               "joshcom.net",
		port:               70,
		rootDirectory:      "testdata/mygopherhole",
		mapFileName:        "gophermap",
		mimeTypeIgnoreList: []string{"application/"},
	}

	t.Run("process file", func(t *testing.T) {
		buffer, writer := newBytesWriter()
		path := "testdata/mygopherhole/art/laptop.txt"
		res, _ := newResource(path)
		_, err := resp.write(writer, &res)

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
		path := "testdata/mygopherhole/phlog/"
		res, _ := newResource(path)
		_, err := resp.write(writer, &res)

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
		path := "testdata/mygopherhole/files"
		res, _ := newResource(path)
		_, err := resp.write(writer, &res)

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

	t.Run("do not process restricted file", func(t *testing.T) {
		_, writer := newBytesWriter()
		path := "testdata/mygopherhole/files/run.exe"
		res, _ := newResource(path)
		_, err := resp.write(writer, &res)

		if err == nil {
			t.Error("Expected error on processing restricted file.")
		}
	})

	t.Run("resource error", func(t *testing.T) {
		buffer, writer := newBytesWriter()
		path := "testdata/mygopherhole/art/iheartsocialmedia.txt"
		res, _ := newResource(path)
		_, err := resp.write(writer, &res)

		if err != nil {
			t.Error("Invalid file should be handled without error.")
		}

		dataStr := buffer.String()
		if strings.Index(dataStr, "3File not found.") < 0 {
			t.Error("Error payload not as expected.")
		}
	})
}
