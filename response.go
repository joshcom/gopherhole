package gopherhole

type response struct {
	host               string
	port               int
	rootDirectory      string
	mapFileName        string
	defaultMime        string
	mimeTypeIgnoreList []string
}

func (r *response) build(res *resource) (pay *[]byte, err error) {
	payloadRes := *res
	var p payload

	if res.isError {
		p = newErrorPayload()
	} else if res.isDirectory {
		dirResources, err := res.directoryResources()
		if err != nil {
			return pay, err
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

	pay, err = p.build(&payloadRes)
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
