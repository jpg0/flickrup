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
	realDateTaken time.Time
	//realDateTaken time.Time
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

func (ftf CachedTaggedFile) ReplaceStringTag(old string, new string) error {

	for i, word := range ftf.stringtags {
		if word == old {
			ftf.stringtags[i] = new
			return nil//only the first
		}
	}

	return nil
}

func (ftf CachedTaggedFile) Keywords() processing.Keywords {
	return ftf.keywords
}

func (ftf CachedTaggedFile) RealDateTaken() time.Time {
	return ftf.realDateTaken
}

func NewFakeTaggedFile(name string, filepath string, stringtags map[string]string, keywords []string, realDateTaken time.Time) *CachedTaggedFile {
	return &CachedTaggedFile{
		name: name,
		filepath: filepath,
		stringtags: stringtags,
		keywords: &CachedKeywords{keywords:keywords},
		realDateTaken: realDateTaken,
	}
}