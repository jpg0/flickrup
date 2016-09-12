package tags

import (
	"testing"
	"github.com/jpg0/flickrup/processing"
	"github.com/jpg0/flickrup/config"
	"github.com/jpg0/flickrup/testlib"
	"time"
	log "github.com/Sirupsen/logrus"
)

func TestNoReplacements(t *testing.T) {

	ctx := applyToTags(map[string]string {
		"a":"b",
		"x":"d",
	}, nil)

	assertEquals("b", ctx.File.StringTag("a"), t)
	assertEquals("d", ctx.File.StringTag("x"), t)
}

func TestDynamicReplacements(t *testing.T) {

	log.SetLevel(log.DebugLevel)

	ctx := applyToTags(map[string]string {
		"y":"z",
	}, []string{"1"})

	assertEquals("2", ctx.File.StringTag("y"), t)
}

func applyToTags(stringtags map[string]string, keywords []string) *processing.ProcessingContext {

	config := &config.Config{
		TagReplacements: map[string]map[string]string{
			"y": map[string]string {"$1": "2"}, //dynamic
		},
	}
	file := testlib.NewFakeTaggedFile("", "", stringtags, keywords, time.Time{})

	ctx := processing.NewProcessingContext(config, file)

	MaybeReplace(ctx)

	return ctx
}

func assertEquals(expected string, actual string, t *testing.T) {
	if actual != expected {
		t.Errorf("Test failed, expected: '%s', got:  '%s'", expected, actual)
	}
}