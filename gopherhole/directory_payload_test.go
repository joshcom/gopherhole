package gopherhole

import (
	"strings"
	"testing"
)

func TestDirectoryPayload_newDirectoryPayload(t *testing.T) {
	payload := newDirectoryPayload("joshcom.net",
		70,
		"testdata/mygopherhole",
		"text/plain",
		[]string{})

	if payload.host != "joshcom.net" {
		t.Errorf("Unexpected host of %s", payload.host)
	}

	if payload.port != 70 {
		t.Errorf("Unexpected port of %d", payload.port)
	}

	if payload.rootDirectory != "testdata/mygopherhole" {
		t.Errorf("Unexpected root directory of %s", payload.rootDirectory)
	}
}

func TestDirectoryPayload_build(t *testing.T) {
	payload := newDirectoryPayload("joshcom.net",
		70,
		"testdata/mygopherhole",
		"text/plain",
		[]string{})

	path := "testdata/mygopherhole/art"
	resource, err := newResource(path)
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	data, err := payload.build(&resource)
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	dataStr := string(*data)

	t.Run("order directory response", func(t *testing.T) {
		fileOrder := []string{
			"cactus.txt",
			"coffee.txt",
			"empty.txt",
			"flipphone.txt",
			"laptop.txt",
			"subdirectory",
		}

		for ind, file := range fileOrder {
			if ind == len(fileOrder)-1 {
				break
			}

			if strings.Index(dataStr, file) > strings.Index(dataStr, fileOrder[ind+1]) {
				t.Errorf("Ordering: %s expected to come before %s.",
					file,
					fileOrder[ind+1])
			}
		}
	})

	t.Run("payload has suffix", func(t *testing.T) {
		if !strings.HasSuffix(dataStr, string(payload.suffix())) {
			t.Error("Expected terminating suffix in payload.")
		}
	})

	t.Run("payload line format", func(t *testing.T) {
		formatedLine := "0cactus.txt\t/art/cactus.txt\tjoshcom.net\t70\r\n"
		if strings.Index(dataStr, formatedLine) == -1 {
			t.Error("Unexpected selector format")
		}
	})

	t.Run("apply default mime types", func(t *testing.T) {
		path := "testdata/mygopherhole/phlog"
		resource, err := newResource(path)

		if err != nil {
			t.Fatalf("Unexpected error %v", err)
		}

		data, err := payload.build(&resource)
		if err != nil {
			t.Fatalf("Unexpected error %v", err)
		}

		dataStr := string(*data)
		if strings.Index(dataStr, "20200908-Its-A-Post.phlog") < 0 {
			t.Error("Expected default mime type to be applied.")
		}
	})

	t.Run("error if building file", func(t *testing.T) {
		path := "testdata/mygopherhole/art/laptop.txt"
		resource, err := newResource(path)
		if err != nil {
			t.Fatalf("Unexpected error %v", err)
		}

		_, err = payload.build(&resource)
		if err == nil {
			t.Error("Expected error when processing file.")
		}
	})

	t.Run("consider blacklist", func(t *testing.T) {
		path := "testdata/mygopherhole/files"
		resource, err := newResource(path)
		if err != nil {
			t.Fatalf("Unexpected error %v", err)
		}

		data, err = payload.build(&resource)
		dataStr := string(*data)
		if err != nil {
			t.Fatalf("Error when building payload")
		}

		gifExpected := "ghappydance.gif"
		htmlExpected := "hindex.html"
		binExpected := "5run.exe"

		for _, expected := range []string{gifExpected, htmlExpected, binExpected} {
			if strings.Index(dataStr, expected) < 0 {
				t.Errorf("Expected to find '%s' in payload", expected)
			}
		}

		payload.mimeTypeIgnoreList = []string{"text/html", "application/"}
		data, err = payload.build(&resource)
		dataStr = string(*data)

		if strings.Index(dataStr, gifExpected) < 0 {
			t.Errorf("Expected to find '%s' in payload", gifExpected)
		}

		for _, expected := range []string{htmlExpected, binExpected} {
			if strings.Index(dataStr, expected) >= 0 {
				t.Errorf("Unexpectedly found '%s' in payload", expected)
			}
		}

	})

	t.Run("build empty directory", func(t *testing.T) {
		payload.mimeTypeIgnoreList = []string{"text/plain"}

		path := "testdata/mygopherhole/art/subdirectory"
		resource, err := newResource(path)
		if err != nil {
			t.Fatalf("Unexpected error %v", err)
		}

		data, err := payload.build(&resource)
		if err != nil {
			t.Fatalf("Unexpected error %v", err)
		}

		dataStr := string(*data)
		if dataStr != string(payload.suffix()) {
			t.Error("Expected payload to be only terminating prefix.")
		}
	})
}
