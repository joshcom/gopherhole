package gopherhole

import (
	"bytes"
)

type mapfilePayload struct {
	*payloadImpl
}

func newMapfilePayload(host string, port int, rootdir string) *mapfilePayload {
	p := mapfilePayload{
		&payloadImpl{
			host:          host,
			port:          port,
			rootDirectory: rootdir,
		},
	}

	return &p
}

func (f *mapfilePayload) build(r *resource) (*payloadReader, error) {
	fileData, err := r.readFileData()
	if err != nil {
		return nil, err
	}

	lines := bytes.Split(*fileData, []byte{LF_CHAR})
	var newLines []byte
	for _, line := range lines {
		line := bytes.TrimSpace(line)
		if !f.isEntityRow(&line) {
			line = *f.buildInlineTextEntityRow(line)
		} else {
			line = *f.correctEntityRow(line, f.host, f.port)
		}
		newLines = append(newLines, line...)
	}

	res := f.pack(&newLines)

	var reader payloadReader = newPayloadBytesReader(res)
	return &reader, err
}
