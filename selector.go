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
	if strings.TrimSpace(s.rootPath) == "" {
		return fullPath, errors.New("Root path cannot be blank.")
	}

	absRoot, err := filepath.Abs(s.rootPath)
	if err != nil {
		return
	}

	fullPath = absRoot + "/" + s.query
	fullPath, err = filepath.Abs(fullPath)
	if err != nil {
		return
	}

	if strings.Index(fullPath, absRoot) != 0 {
		err = fmt.Errorf("Evaluated path %s does not begin with specified root %s",
			fullPath,
			absRoot)
		return
	}

	return
}
