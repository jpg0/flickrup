package processing

import "github.com/juju/errors"

type TaggedFileFactory struct {
	Constructors map[string]func (filepath string) (TaggedFile, error)
}

func (tff TaggedFileFactory) LoadTaggedFile(path string) (TaggedFile, error) {
	constructor := tff.Constructors[path]

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