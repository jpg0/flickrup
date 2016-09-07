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

	ctx := processing.NewProcessingContext()
	ctx.Config = &config.Config{}
	ctx.File = testlib.NewFakeTaggedFile("", "", nil, []string{"sharing:visibility=private"}, time.Time{})

	NewRewriter().MaybeRewrite(ctx)

	expected := []string{"sharing:visibility=private"}
	actual := ctx.File.Keywords().All()
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Test failed, expected: '%s', got:  '%s'", expected, actual)
	}
}

func TestRewriteRequired(t *testing.T){

	ctx := processing.NewProcessingContext()

	ctx.File = testlib.NewFakeTaggedFile("", "", nil, []string{"sharing:visibility::private"}, time.Time{})

	NewRewriter().MaybeRewrite(ctx)

	expected := []string{"sharing:visibility=private"}
	actual := ctx.File.Keywords().All()
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Test failed, expected: '%s', got:  '%s'", expected, actual)
	}
}

