package testlib

import (
	"time"
	"github.com/jpg0/flickrup/processing"
)

type CachedTaggedFile struct {
	name string
	filepath string
	stringtags map[string]string
	keywords *CachedKeywords
	dateTaken time.Time
}

type CachedKeywords struct {
	keywords []string
}

func (fk CachedKeywords) All() []string {
	return fk.keywords
}

func (fk CachedKeywords) Replace(old string, new string) error {
	for i, word := range fk.keywords {
		if word == old {
			fk.keywords[i] = new
			return nil//only the first
		}
	}

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

func (ftf CachedTaggedFile) Keywords() processing.Keywords {
	return ftf.keywords
}

func (ftf CachedTaggedFile) DateTaken() time.Time {
	return ftf.dateTaken
}

func NewFakeTaggedFile(name string, filepath string, stringtags map[string]string, keywords []string, dateTaken time.Time) *CachedTaggedFile {
	return &CachedTaggedFile{
		name: name,
		filepath: filepath,
		stringtags: stringtags,
		keywords: &CachedKeywords{keywords:keywords},
		dateTaken: dateTaken,
	}
}