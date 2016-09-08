package tags

import (
	"testing"
	"github.com/jpg0/flickrup/processing"
	"github.com/jpg0/flickrup/config"
	"github.com/jpg0/flickrup/testlib"
	"time"
)

func applyToTags(stringtags map[string]string) *processing.ProcessingContext {
	ctx := processing.NewProcessingContext()
	ctx.Config = &config.Config{
		TagReplacements: map[string]map[string]string{
			"x": map[string]string {"1": "2"},
			"y": map[string]string {"3": "4"},
		},
	}
	ctx.File = testlib.NewFakeTaggedFile("", "", stringtags, nil, time.Time{})

	MaybeReplace(ctx)

	return ctx
}

func TestNoReplacements(t *testing.T) {


	ctx := applyToTags(map[string]string {
		"a":"b",
		"c":"d",
	})

	assertEquals("b", ctx.File.StringTag("a"), t)
	assertEquals("d", ctx.File.StringTag("c"), t)
}

func assertEquals(expected string, actual string, t *testing.T) {
	if actual != expected {
		t.Errorf("Test failed, expected: '%s', got:  '%s'", expected, actual)
	}
}