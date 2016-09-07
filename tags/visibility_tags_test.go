package tags

import (
	"testing"
	"github.com/jpg0/flickrup/processing"
	"github.com/jpg0/flickrup/config"
	"github.com/jpg0/flickrup/testlib"
	"time"
)

func TestNoVisibilityRequired(t *testing.T){

	ctx := processing.NewProcessingContext()
	ctx.Config = &config.Config{}

	ExtractVisibility(ctx)

	expected := "public"
	actual := ctx.Visibilty
	if actual != expected {
		t.Errorf("Test failed, expected: '%s', got:  '%s'", expected, actual)
	}
}

func TestPrivateVisibilityRequired(t *testing.T){

	ctx := processing.NewProcessingContext()
	ctx.Config = &config.Config{
		VisibilityPrefix: "sharing:visibility=",
	}

	ctx.File = testlib.NewFakeTaggedFile("", "", nil, []string{"sharing:visibility=private"}, time.Time{})

	ExtractVisibility(ctx)

	expected := "private"
	actual := ctx.Visibilty
	if actual != expected {
		t.Errorf("Test failed, expected: '%s', got:  '%s'", expected, actual)
	}
}

