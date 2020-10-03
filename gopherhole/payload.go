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

func (f *payloadImpl) correctEntityRow(row []byte, host string, port int) *[]byte {
	tabColumns := bytes.Split(row, []byte{TAB_CHAR})

	// If this row already has 4 or more columns, it's as correct as its
	// gonna get.
	// If there's not at least a name and a path, we won't apply the
	// local host and port, and leave luck to heaven.
	if len(tabColumns) < 2 || len(tabColumns) >= 4 {
		return &row
	}

	// If this row already has a host specified, use that.
	var hostCol []byte
	if len(tabColumns) > 2 {
		hostCol = bytes.TrimSpace(tabColumns[2])
	} else {
		hostCol = []byte(host)
	}

	name := bytes.TrimSpace(tabColumns[0])
	path := bytes.TrimSpace(tabColumns[1])

	return f.buildRow(&name, &path, &hostCol, port)
}

func (f *payloadImpl) buildEntityRow(entity EntityType, name string, path string, host string, port int) *[]byte {
	entityName := []byte{byte(entity)}
	entityName = append(entityName, []byte(name)...)

	pathCol := []byte(path)
	hostCol := []byte(host)
	return f.buildRow(&entityName, &pathCol, &hostCol, port)
}

func (f *payloadImpl) buildRow(entityName *[]byte, path *[]byte, host *[]byte, port int) *[]byte {
	p := fmt.Sprintf("%d", port)
	portCol := []byte(p)

	row := []byte{}
	row = append(row, *entityName...)
	row = append(row, TAB_CHAR)
	row = append(row, *path...)
	row = append(row, TAB_CHAR)
	row = append(row, *host...)
	row = append(row, TAB_CHAR)
	row = append(row, portCol...)
	row = append(row, CR_CHAR, LF_CHAR)
	return &row
}
