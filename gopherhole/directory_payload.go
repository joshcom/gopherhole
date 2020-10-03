package gopherhole

import (
	"sort"
)

type directoryPayload struct {
	*payloadImpl
}

func newDirectoryPayload(host string,
	port int,
	rootdir string,
	defaultMime string,
	mimeTypeIgnoreList []string) *directoryPayload {

	p := directoryPayload{
		&payloadImpl{
			host:               host,
			port:               port,
			rootDirectory:      rootdir,
			defaultMime:        defaultMime,
			mimeTypeIgnoreList: mimeTypeIgnoreList,
		},
	}

	return &p
}

func (f *directoryPayload) build(r *resource) (*[]byte, error) {
	res := new([]byte)
	directoryResources, err := r.directoryResources()
	if err != nil {
		return res, err
	}

	sort.Slice(directoryResources, func(a, b int) bool {
		return directoryResources[a].name < directoryResources[b].name
	})

	for _, resource := range directoryResources {
		if resource.mimeType.in(f.mimeTypeIgnoreList) {
			continue
		}
		row := f.buildResourceEntityRow(&resource)
		*res = append(*res, (*row)...)
	}

	res = f.pack(res)
	return res, err
}
