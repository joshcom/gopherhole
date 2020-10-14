package gopherhole

import (
	"errors"
	"os"
)

// Implements payload reader
type payloadFileReader struct {
	reader *os.File
}

func newPayloadFileReader(file *os.File) *payloadFileReader {
	return &payloadFileReader{
		reader: file,
	}
}

func (r *payloadFileReader) Close() (err error) {
	return r.reader.Close()
}

func (r *payloadFileReader) Read(p []byte) (n int, err error) {
	return r.reader.Read(p)
}

type filePayload struct {
	*payloadImpl
}

func newFilePayload(defaultMime string, mimeTypeIgnoreList []string) *filePayload {
	pay := filePayload{
		&payloadImpl{
			defaultMime:        defaultMime,
			mimeTypeIgnoreList: mimeTypeIgnoreList,
		},
	}
	return &pay
}

func (f *filePayload) build(r *resource) (*payloadReader, error) {
	if r.mimeType.in(f.mimeTypeIgnoreList) {
		return nil, errors.New("File access is forbidden due to mime type restrictions.")
	}

	file, err := r.file()
	if err != nil {
		return nil, err
	}

	var reader payloadReader = newPayloadFileReader(file)
	return &reader, err
}
