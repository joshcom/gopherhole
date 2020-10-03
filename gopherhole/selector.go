package gopherhole

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
)

type selector struct {
	query    string
	rootPath string
}

func newSelector(rootPath string, query string) selector {
	query = strings.TrimSpace(query)
	selector := selector{rootPath: rootPath, query: query}
	return selector
}

func (s *selector) path() (fullPath string, err error) {
	// In the future, this check ought to be made during a config
	// validation check.
	if strings.TrimSpace(s.rootPath) == "" {
		return fullPath, errors.New("Root path cannot be blank.")
	}

	// Make an absolute path out of the root path configured for
	// this server.  This will change relative paths to absolute
	// paths, for the check that comes up next.
	absRoot, err := filepath.Abs(s.rootPath)
	if err != nil {
		return
	}

	// Make a full, and absolute, path out of the root path
	// + the query ('/var/gopher' + '/phlog' = '/var/gopher/phlog').
	fullPath = absRoot + "/" + s.query
	fullPath, err = filepath.Abs(fullPath)
	if err != nil {
		return
	}

	// Finally, ensure the derived path to the file is within the
	// root directory.
	if strings.Index(fullPath, absRoot) != 0 {
		err = fmt.Errorf("Evaluated path %s does not begin with specified root %s",
			fullPath,
			absRoot)
		return
	}

	return
}
