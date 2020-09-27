package gopherhole

import (
	"strings"
	"testing"
)

func TestFilePayload_BuildResponse(t *testing.T) {
	payload := new(filePayload)

	t.Run("build response for file", func(t *testing.T) {
		path := "testdata/mygopherhole/gophermap"
		resource, err := newResource(path)
		if err != nil {
			t.Fatalf("Unexpected error %v", err)
		}

		data, err := payload.build(&resource)
		if err != nil {
			t.Fatalf("Unexpected error %v", err)
		}

		dataStr := string(*data)
		if strings.Index(dataStr, "about.txt") < 0 {
			t.Error("Payload data not as expected.")
		}

		if strings.HasSuffix(dataStr, string(payload.suffix())) {
			t.Error("Terminating suffix should not be applied to files.")
		}
	})

	t.Run("build response for empty file", func(t *testing.T) {
		path := "testdata/mygopherhole/art/empty.txt"
		resource, err := newResource(path)
		if err != nil {
			t.Fatalf("Unexpected error %v", err)
		}

		data, err := payload.build(&resource)
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
		resource, err := newResource(path)
		if err != nil {
			t.Fatalf("Unexpected error %v", err)
		}

		_, err = payload.build(&resource)
		if err == nil {
			t.Error("Expected error when processing directory.")
		}
	})
}
