package gopherhole

import (
	"bytes"
	"io"
	"io/ioutil"
)

func buildConfiguration() (config Configuration) {
	config = NewConfiguration()
	config.RootDirectory = "testdata/mygopherhole"
	return config
}

func readAllString(reader io.Reader) string {
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		panic(err)
	}
	return string(data)
}

func newBytesWriter() (*bytes.Buffer, *io.Writer) {
	var b bytes.Buffer
	writer := io.Writer(&b)
	return &b, &writer
}
