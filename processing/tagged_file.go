package processing

import (
	"time"
	"strings"
)

type TaggedFile interface {
	Filepath() string
	Name() string
	Keywords() Keywords
	RealDateTaken() time.Time
	StringTag(name string) string
	ReplaceStringTag(old string, new string) error
}

type Keywords interface {
	All() []string
	Replace(old string, new string) error
}

type KeywordsHelper struct {
	Keywords
}

func ValuesByPrefix(k Keywords, prefix string) []string {
	values := make([]string, 0)
	for _, v := range k.All() {
		if strings.HasPrefix(v, prefix) {
			values = append(values, v[len(prefix):])
		}
	}
	return values
}



//type TaggedVideo struct {
//	file *os.File
//	tags []string
//}
//
//func (ti TaggedVideo) File() *os.File {
//	return ti.file
//}
//
//func (ti TaggedVideo) Keywords() []string {
//	return ti.tags
//}
