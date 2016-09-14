package processing

import (
	"path/filepath"
	"strings"
	"fmt"
)

type TaggedFileFactory struct {
	Constructors map[string]func (filepath string) (TaggedFile, error)
}

type NoConstructorAvailableError struct {
	path string
}

func (ncae NoConstructorAvailableError) Error() string {
	return fmt.Sprintf("No constructor for file: %v", ncae.path)
}

func (tff TaggedFileFactory) LoadTaggedFile(path string) (TaggedFile, error) {
	ext := strings.ToLower(filepath.Ext(path))

	if len(ext) == 0 {
		return nil, NoConstructorAvailableError{path}
	}

	constructor := tff.Constructors[ext[1:]]

	if constructor == nil {
		return nil, NoConstructorAvailableError{path}
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