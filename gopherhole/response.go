package gopherhole

import (
	"io"
)

type response struct {
	host               string
	port               int
	rootDirectory      string
	mapFileName        string
	defaultMime        string
	mimeTypeIgnoreList []string
}

func (r *response) write(conn *io.Writer, res *resource) (written int64, err error) {
	payloadRes := *res
	var p payload

	if res.isError {
		p = newErrorPayload()
	} else if res.isDirectory {
		dirResources, err := res.directoryResources()
		if err != nil {
			return written, err
		}
		mapfile := r.findMapFile(&dirResources)
		if mapfile != nil {
			p = newMapfilePayload(
				r.host,
				r.port,
				r.rootDirectory,
			)
			payloadRes = *mapfile
		} else {
			p = newDirectoryPayload(
				r.host,
				r.port,
				r.rootDirectory,
				r.defaultMime,
				r.mimeTypeIgnoreList,
			)
		}
	} else {
		p = newFilePayload(r.defaultMime, r.mimeTypeIgnoreList)
	}

	reader, err := p.build(&payloadRes)
	if err != nil {
		return
	}

	defer (*reader).Close()
	written, err = io.Copy(*conn, *reader)

	return
}

func (r *response) findMapFile(resList *[]resource) (mapfile *resource) {
	for _, res := range *resList {
		if res.name == r.mapFileName {
			mapfile = &res
			break
		}
	}
	return
}
