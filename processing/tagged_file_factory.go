package processing

import (
	"github.com/juju/errors"
	"path/filepath"
	"strings"
)

type TaggedFileFactory struct {
	Constructors map[string]func (filepath string) (TaggedFile, error)
}

func (tff TaggedFileFactory) LoadTaggedFile(path string) (TaggedFile, error) {
	ext := strings.ToLower(filepath.Ext(path))

	if len(ext) == 0 {
		return nil, errors.Errorf("No constructor for file: %v", path)
	}

	constructor := tff.Constructors[ext[1:]]

	if constructor == nil {
		return nil, errors.Errorf("No constructor for file: %v", path)
	}

	return constructor(path)
}

func MergeTaggedFileFactories(factories ...*TaggedFileFactory) *TaggedFileFactory {
	constructors := make(map[string]func (filepath string) (TaggedFile, error))

	for _, f := range factories {
		for k, v := range f.Constructors {
			constructors[k] = v
		}
	}

	return &TaggedFileFactory{
		Constructors: constructors,
	}
}