package gopherhole

import (
	"strings"
	"testing"
)

func TestFilePayload_BuildResponse(t *testing.T) {
	pay := newFilePayload("text/plain", []string{"text/html"})

	t.Run("build response for file", func(t *testing.T) {
		path := "testdata/mygopherhole/gophermap"
		res, err := newResource(path)
		if err != nil {
			t.Fatalf("Unexpected error %v", err)
		}

		data, err := pay.build(&res)
		if err != nil {
			t.Fatalf("Unexpected error %v", err)
		}

		dataStr := string(*data)
		if strings.Index(dataStr, "about.txt") < 0 {
			t.Error("Payload data not as expected.")
		}

		if strings.HasSuffix(dataStr, string(pay.suffix())) {
			t.Error("Terminating suffix should not be applied to files.")
		}
	})

	t.Run("build response for empty file", func(t *testing.T) {
		path := "testdata/mygopherhole/art/empty.txt"
		res, err := newResource(path)
		if err != nil {
			t.Fatalf("Unexpected error %v", err)
		}

		data, err := pay.build(&res)
		if err != nil {
			t.Fatalf("Unexpected error %v", err)
		}

		dataStr := string(*data)
		if dataStr != "" {
			t.Error("Expected empty data string.")
		}
	})

	t.Run("error when building response for directory", func(t *testing.T) {
		path := "testdata/mygopherhole/art"
		res, err := newResource(path)
		if err != nil {
			t.Fatalf("Unexpected error %v", err)
		}

		_, err = pay.build(&res)
		if err == nil {
			t.Error("Expected error when processing directory.")
		}
	})

	t.Run("error when building response for file of forbidden mime type", func(t *testing.T) {
		path := "testdata/mygopherhole/files/index.html"
		res, _ := newResource(path)

		_, err := pay.build(&res)
		if err == nil {
			t.Error("Expected error accessing forbidden file.")
		}
	})
}
