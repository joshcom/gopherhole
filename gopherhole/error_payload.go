package gopherhole

type errorPayload struct {
	*payloadImpl
}

func newErrorPayload() *errorPayload {
	p := errorPayload{}
	return &p
}

func (f *errorPayload) build(r *resource) (*payloadReader, error) {
	row := f.buildResourceEntityRow(r)
	res := f.pack(row)

	var reader payloadReader = newPayloadBytesReader(res)
	return &reader, nil
}
