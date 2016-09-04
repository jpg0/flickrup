package filetype

import (
)

type TaggedFile interface {
	Filepath() string
	Name() string
	Keywords() []string
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
