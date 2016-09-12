package testlib

import (
	"time"
	"github.com/jpg0/flickrup/processing"
	"fmt"
)

type CachedTaggedFile struct {
	name string
	filepath string
	stringtags map[string]string
	keywords *CachedKeywords
	realDateTaken time.Time
	//realDateTaken time.Time
}

type CachedKeywords struct {
	keywords *processing.TagSet
}

func (fk CachedKeywords) All() *processing.TagSet {
	return fk.keywords
}

func (fk CachedKeywords) Replace(old string, new string) error {

	fk.keywords.Remove(old)
	fk.keywords.Add(new)

	return nil
}

func (ftf CachedTaggedFile) Name() string {
	return ftf.name
}

func (ftf CachedTaggedFile) Filepath() string {
	return ftf.filepath
}

func (ftf CachedTaggedFile) StringTag(name string) string {
	return ftf.stringtags[name]
}

func (ftf CachedTaggedFile) ReplaceStringTag(name string, new string) error {

	ftf.stringtags[name] = new

	return nil
}

func (ftf CachedTaggedFile) Keywords() processing.Keywords {
	return ftf.keywords
}

func (ftf CachedTaggedFile) DateTaken() time.Time {
	return ftf.realDateTaken
}

func NewFakeTaggedFile(name string, filepath string, stringtags map[string]string, keywords []string, realDateTaken time.Time) *CachedTaggedFile {
	return &CachedTaggedFile{
		name: name,
		filepath: filepath,
		stringtags: stringtags,
		keywords: &CachedKeywords{processing.NewTagSet(keywords)},
		realDateTaken: realDateTaken,
	}
}

func Dump(file processing.TaggedFile) {
	cf := file.(*CachedTaggedFile)

	fmt.Println(cf.stringtags)
}