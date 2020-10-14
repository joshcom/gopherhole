package gopherhole

import (
	"io"
)

type request struct {
	configuration Configuration
	selector      selector
}

func newRequest(query string, config Configuration) (req request) {
	selector := newSelector(config.RootDirectory, query)

	req = request{
		configuration: config,
		selector:      selector,
	}
	return req
}

func (r *request) process(conn *io.Writer) (written int64, err error) {
	fullPath, err := r.selector.path()
	var res resource

	// Errors at this state will be a problem
	// with the selector query itself.
	if err != nil {
		res, _ = newQueryErrorResource(r.selector.query)
	} else {
		res, _ = newResource(fullPath)
	}

	written, err = r.writePayload(conn, &res)

	// Errors here mean the selector was constructed
	// fine, but the file itself is forbidden.
	if err != nil {
		res, _ = newNotFoundErrorResource(r.selector.query)
		written, err = r.writePayload(conn, &res)
	}

	return
}

func (r *request) writePayload(conn *io.Writer, res *resource) (written int64, err error) {
	conf := r.configuration
	resp := response{
		host:               conf.Host,
		port:               conf.Port,
		rootDirectory:      conf.RootDirectory,
		mapFileName:        conf.MapFileName,
		defaultMime:        conf.DefaultMimeType,
		mimeTypeIgnoreList: conf.MimeTypeIgnoreList,
	}

	return resp.write(conn, res)
}
