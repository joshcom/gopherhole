package gopherhole

import (
	"bytes"
	"fmt"
	"strings"
)

const (
	CR_CHAR  = '\r'
	END_CHAR = '.'
	LF_CHAR  = '\n'
	TAB_CHAR = '\t'
	NOPATH   = "(NOTHING)"
	NOHOST   = "nohost"
	NOPORT   = 0
)

type payload interface {
	build(*resource) (*[]byte, error)
	pack(*[]byte) *[]byte
	suffix() []byte
}

type payloadImpl struct {
	rootDirectory      string
	host               string
	port               int
	defaultMime        string
	mimeTypeIgnoreList []string
}

func (r *payloadImpl) build(resource *resource) (p *[]byte, err error) {
	return
}

func (r *payloadImpl) pack(p *[]byte) *[]byte {
	response := *p
	response = append(response, r.suffix()...)
	return &response
}

func (r *payloadImpl) suffix() (suffix []byte) {
	suffix = append(suffix, CR_CHAR, LF_CHAR, END_CHAR, CR_CHAR, LF_CHAR)
	return suffix
}

// We'll keep it simple.  If there's a tab at all, it's an
// entity row.
func (f *payloadImpl) isEntityRow(row *[]byte) bool {
	tabColumns := bytes.Split(*row, []byte{TAB_CHAR})
	if len(tabColumns) <= 1 {
		return false
	} else {
		return true
	}
}

func (f *payloadImpl) buildInlineTextEntityRow(row []byte) *[]byte {
	cleansedRow := string(row)
	cleansedRow = strings.TrimSpace(cleansedRow)
	newRow := f.buildEntityRow(InlineTextEntity, cleansedRow, NOPATH, NOHOST, NOPORT)
	return newRow
}

func (f *payloadImpl) buildResourceEntityRow(res *resource) *[]byte {
	row := new([]byte)
	var entity EntityType
	ok := true

	if res.isError {
		return f.buildErrorEntityRow(res.name)
	} else if res.isDirectory {
		entity = DirectoryEntity
	} else {
		entity, ok = EntityForMimeType(res.mimeType.string())
		if !ok {
			entity, ok = EntityForMimeType(f.defaultMime)
		}
	}

	if !ok {
		return row
	}

	path := res.path
	path = strings.Replace(path, f.rootDirectory, "", 1)

	row = f.buildEntityRow(entity, res.name, path, f.host, f.port)

	return row
}

func (f *payloadImpl) buildErrorEntityRow(message string) *[]byte {
	row := f.buildEntityRow(ErrorEntity, message, NOPATH, NOHOST, NOPORT)
	return row
}

func (f *payloadImpl) buildEntityRow(entity EntityType, name string, path string, host string, port int) *[]byte {
	portCol := fmt.Sprintf("%d", port)

	row := []byte{byte(entity)}
	row = append(row, []byte(name)...)
	row = append(row, TAB_CHAR)
	row = append(row, []byte(path)...)
	row = append(row, TAB_CHAR)
	row = append(row, []byte(host)...)
	row = append(row, TAB_CHAR)
	row = append(row, []byte(portCol)...)
	row = append(row, CR_CHAR, LF_CHAR)
	return &row
}
