package tags

import (
	"testing"
	"github.com/jpg0/flickrup/processing"
	"github.com/jpg0/flickrup/config"
	"github.com/jpg0/flickrup/testlib"
	"time"
	"reflect"
)

func TestNoRewriteRequired(t *testing.T){

	config := &config.Config{}
	file := testlib.NewFakeTaggedFile("", "", nil, []string{"sharing:visibility=private"}, time.Time{})

	ctx := processing.NewProcessingContext(config, file)


	NewRewriter().MaybeRewrite(ctx)

	expected := []string{"sharing:visibility=private"}
	actual := ctx.File.Keywords().All().Slice()
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Test failed, expected: '%s', got:  '%s'", expected, actual)
	}
}

func TestRewriteRequired(t *testing.T){

	config := &config.Config{}
	file := testlib.NewFakeTaggedFile("", "", nil, []string{"sharing:visibility::private"}, time.Time{})

	ctx := processing.NewProcessingContext(config, file)

	NewRewriter().MaybeRewrite(ctx)

	expected := []string{"sharing:visibility=private"}
	actual := ctx.File.Keywords().All().Slice()
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Test failed, expected: '%s', got:  '%s'", expected, actual)
	}
}

