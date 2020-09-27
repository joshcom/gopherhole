package gopherhole

import (
	"errors"
)

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

func (f *filePayload) build(r *resource) (pay *[]byte, err error) {
	if r.mimeType.in(f.mimeTypeIgnoreList) {
		return nil, errors.New("File access is forbidden due to mime type restrictions.")
	}

	pay, err = r.readFileData()
	if err != nil {
		return
	}
	return
}
