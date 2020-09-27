package gopherhole

type filePayload struct {
	*payloadImpl
}

func newFilePayload() *filePayload {
	payload := filePayload{}
	return &payload
}

func (f *filePayload) build(r *resource) (payload *[]byte, err error) {
	payload, err = r.readFileData()
	if err != nil {
		return
	}
	return
}
