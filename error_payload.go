package gopherhole

type errorPayload struct {
	*payloadImpl
}

func newErrorPayload() *errorPayload {
	p := errorPayload{}
	return &p
}

func (f *errorPayload) build(r *resource) (res *[]byte, err error) {
	row := f.buildResourceEntityRow(r)
	res = f.pack(row)
	return
}
