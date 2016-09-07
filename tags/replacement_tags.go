package tags

import (
	"github.com/jpg0/flickrup/processing"
	"strings"
)

func MaybeReplace(ctx *processing.ProcessingContext) error {
	replacements := ctx.Config.TagReplacements
	allKeywords := ctx.File.Keywords().All()

	if replacements != nil {
		for tagName, tagValueReplacements := range replacements {
			existing := ctx.File.StringTag(tagName)

			for key, value := range tagValueReplacements {
				if strings.HasPrefix(key, "$") {
					if contains(allKeywords, key[1:]) {
						err := ctx.File.ReplaceStringTag(tagName, value)

						if err != nil {
							return err
						}
					}
				} else if existing == key {
					err := ctx.File.ReplaceStringTag(tagName, value)

					if err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}