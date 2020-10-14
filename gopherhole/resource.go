package gopherhole

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type resource struct {
	path        string
	name        string
	mimeType    *mimeType
	isDirectory bool
	isError     bool
}

func newQueryErrorResource(query string) (res resource, err error) {
	message := fmt.Sprintf("Invalid query: %s", query)
	res = resource{
		name:        message,
		isDirectory: false,
		isError:     true,
	}

	return
}

func newNotFoundErrorResource(query string) (res resource, err error) {
	message := fmt.Sprintf("File not found. (%s)", query)
	res = resource{
		name:        message,
		isDirectory: false,
		isError:     true,
	}

	return
}

func newResource(absPath string) (res resource, err error) {
	fileInfo, err := fileStatFromPath(absPath)

	var name string
	var isDirectory, isError bool
	mime := new(mimeType)

	if err == nil {
		name = fileInfo.Name()
		mime.guess(name)
		isDirectory = fileInfo.IsDir()
		isError = false
	} else {
		name = "File not found."
		isError = true
	}

	res = resource{
		path:        absPath,
		name:        name,
		mimeType:    mime,
		isDirectory: isDirectory,
		isError:     isError,
	}

	return
}

func (r *resource) file() (f *os.File, err error) {
	if r.isDirectory {
		return nil, fmt.Errorf("'%s' is a directory.", r.path)
	}

	return os.Open(r.path)
}

func (r *resource) readFileData() (data *[]byte, err error) {
	if r.isDirectory {
		return nil, fmt.Errorf("'%s' is a directory.", r.path)
	}

	data, err = readFileDataFromPath(r.path)
	return
}

func (r *resource) directoryResources() (resources []resource, err error) {
	if !r.isDirectory {
		return nil, fmt.Errorf("'%s' is not a directory.", r.path)
	}

	fileInfos, err := readDirFromPath(r.path)
	if err != nil {
		return
	}

	for _, fileInfo := range fileInfos {
		path := r.path + "/" + fileInfo.Name()
		path = filepath.Clean(path)

		resource, err := newResource(path)
		if err != nil {
			continue
		}

		resources = append(resources, resource)
	}

	return resources, nil
}

var fileStatFromPath = func(absPath string) (fileInfo os.FileInfo, err error) {
	file, err := os.Open(absPath)
	if err != nil {
		return
	}
	defer file.Close()
	fileInfo, err = file.Stat()
	return fileInfo, err
}

var readFileDataFromPath = func(path string) (data *[]byte, err error) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	readData, err := ioutil.ReadAll(file)
	return &readData, err
}

var readDirFromPath = func(absPath string) (fileInfos []os.FileInfo, err error) {
	directory, err := os.Open(absPath)
	if err != nil {
		return
	}
	defer directory.Close()

	fileInfos, err = directory.Readdir(-1)
	return
}
