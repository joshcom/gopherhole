package gopherhole

import (
	"strings"
	"testing"
)

func TestResource_newResource(t *testing.T) {
	t.Run("file resource", func(t *testing.T) {
		path := "testdata/mygopherhole/gophermap"
		res, err := newResource(path)
		if err != nil {
			t.Fatalf("Unexpected error %v", err)
		}

		if res.name != "gophermap" {
			t.Errorf("Unexpected name of %s", res.name)
		}

		if !res.mimeType.isNone() {
			t.Errorf("Unexpected mime type %s.", res.mimeType.string())
		}

		if res.isDirectory {
			t.Error("Resource should not be directory")
		}

		if res.path != path {
			t.Errorf("Expected path of %s, found %s.",
				path,
				res.path)
		}
	})

	t.Run("file with extension resource", func(t *testing.T) {
		path := "testdata/mygopherhole/art/laptop.txt"
		res, _ := newResource(path)

		if res.mimeType.isNone() || res.mimeType.string() != "text/plain" {
			t.Errorf("Unexpected mime type %s.", res.mimeType.string())
		}
	})

	t.Run("directory resource", func(t *testing.T) {
		path := "testdata/mygopherhole"
		res, err := newResource(path)
		if err != nil {
			t.Fatalf("Unexpected error %v", err)
		}

		if res.name != "mygopherhole" {
			t.Errorf("Unexpected name of %s", res.name)
		}

		if !res.mimeType.isNone() {
			t.Errorf("Unexpected mime type %s.", res.mimeType.string())
		}

		if !res.isDirectory {
			t.Error("Resource should be a directory")
		}

		if res.path != path {
			t.Errorf("Expected path of %s, found %s.",
				path,
				res.path)
		}
	})

	t.Run("file not found", func(t *testing.T) {
		res, err := newResource("testdata/mygopherhole/notfound.txt")
		if err == nil {
			t.Error("Expected error for missing file.")
		}

		if !res.isError {
			t.Error("Resource expected to be flagged as error")
		}

		if !res.mimeType.isNone() {
			t.Errorf("Unexpected mime type %s.", res.mimeType.string())
		}
	})
}

func TestResource_newQueryErrorResource(t *testing.T) {
	res, _ := newQueryErrorResource("/blog")

	if !res.isError {
		t.Error("Resource expected to be flagged as error")
	}

	if strings.Index(res.name, "/blog") < 0 {
		t.Error("Expected original query to be part of resource name.")
	}
}

func TestResource_readFileData(t *testing.T) {
	t.Run("read file data", func(t *testing.T) {
		res, err := newResource("testdata/mygopherhole/art/laptop.txt")
		if err != nil {
			t.Fatalf("Unexpected error %v", err)
		}

		data, err := res.readFileData()
		if err != nil {
			t.Fatalf("Unexpected error %v", err)
		}

		if strings.Index(string(*data), "$ vim") < 0 {
			t.Error("File not loaded correctly.")
		}
	})

	t.Run("cannot invoke on directories", func(t *testing.T) {
		res, err := newResource("testdata/mygopherhole/art")
		if err != nil {
			t.Fatalf("Unexpected error %v", err)
		}

		_, err = res.readFileData()
		if err == nil {
			t.Error("Expected error when calling readFileData on directory")
		}
	})
}

func TestResource_directoryResources(t *testing.T) {
	t.Run("list contents of directory", func(t *testing.T) {
		res, _ := newResource("testdata/mygopherhole/art")
		data, _ := res.directoryResources()

		if len(data) != 6 {
			t.Errorf("Expected 5 files in directory, found %d.", len(data))
		}

		var file *resource
		for _, r := range data {
			r := r
			if r.name == "laptop.txt" {
				file = &r
			}
		}

		if file == nil {
			t.Error("Expected to file file 'laptop.txt'")
		}

		if file.path != "testdata/mygopherhole/art/laptop.txt" {
			t.Errorf("Path '%s' not set as expected", file.path)
		}

		if file.mimeType.string() != "text/plain" {
			t.Errorf("Unexpected mime type: '%s'", file.mimeType.string())
		}
	})

	t.Run("list contents of empty directory", func(t *testing.T) {
		resource, _ := newResource("testdata/mygopherhole/empty")
		data, _ := resource.directoryResources()

		if len(data) != 0 {
			t.Error("Expected empty set of resources for empty directory")
		}
	})

	t.Run("cannot invoke unless directory", func(t *testing.T) {
		resource, err := newResource("testdata/mygopherhole/art/laptop.txt")
		if err != nil {
			t.Fatalf("Unexpected error %v", err)
		}

		_, err = resource.directoryResources()
		if err == nil {
			t.Error("Expected error when calling DirectoryResources on regular file")
		}
	})
}
