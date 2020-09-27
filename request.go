package gopherhole

type request struct {
	configuration Configuration
	selector      selector
	payload       *[]byte
}

func newRequest(query string, config Configuration) (req request) {
	selector := newSelector(config.RootDirectory, query)

	req = request{
		configuration: config,
		selector:      selector,
	}
	return req
}

func (r *request) process() (err error) {
	fullPath, err := r.selector.path()
	var res resource

	if err != nil {
		res, _ = newQueryErrorResource(r.selector.query)
	} else {
		res, _ = newResource(fullPath)
	}

	r.payload, err = r.getPayload(&res)

	return
}

func (r *request) getPayload(res *resource) (payload *[]byte, err error) {
	conf := r.configuration
	resp := response{
		host:               conf.Host,
		port:               conf.Port,
		rootDirectory:      conf.RootDirectory,
		mapFileName:        conf.MapFileName,
		defaultMime:        conf.DefaultMimeType,
		mimeTypeIgnoreList: conf.MimeTypeIgnoreList,
	}
	return resp.build(res)
}
