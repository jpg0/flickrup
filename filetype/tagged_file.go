package filetype

import (
	"os"
	"context"
)

type TaggedFile interface {
	File() *os.File
	Keywords() []string
}

type TaggedImage struct {
	file *os.File
	tags []string
}

func (ti TaggedImage) File() *os.File {
	return ti.file
}

func (ti TaggedImage) Keywords() []string {
	return ti.tags
}

type TaggedVideo struct {
	file *os.File
	tags []string
}

func (ti TaggedVideo) File() *os.File {
	return ti.file
}

func (ti TaggedVideo) Keywords() []string {
	return ti.tags
}

func StringValueWithDefault(ctx context.Context, name interface{}, defaultValue string) string {
	rv := ctx.Value(name)

	if rv == nil {
		return defaultValue
	} else {
		return rv.(string)
	}
}

func BoolValueWithDefault(ctx context.Context, name interface{}, defaultValue bool) bool {
	rv := ctx.Value(name)

	if rv == nil {
		return defaultValue
	} else {
		return rv.(bool)
	}
}