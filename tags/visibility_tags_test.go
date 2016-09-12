package tags

import (
	"testing"
	"github.com/jpg0/flickrup/processing"
	"github.com/jpg0/flickrup/config"
	"github.com/jpg0/flickrup/testlib"
	"time"
)

func TestNoVisibilityRequired(t *testing.T){

	ctx := processing.NewProcessingContext(&config.Config{}, nil)

	ExtractVisibility(ctx)

	expected := "public"
	actual := ctx.Visibilty
	if actual != expected {
		t.Errorf("Test failed, expected: '%s', got:  '%s'", expected, actual)
	}
}

func TestPrivateVisibilityRequired(t *testing.T){

	config := &config.Config{
		VisibilityPrefix: "sharing:visibility=",
	}

	file := testlib.NewFakeTaggedFile("", "", nil, []string{"sharing:visibility=private"}, time.Time{})

	ctx := processing.NewProcessingContext(config, file)


	ExtractVisibility(ctx)

	expected := "private"
	actual := ctx.Visibilty
	if actual != expected {
		t.Errorf("Test failed, expected: '%s', got:  '%s'", expected, actual)
	}
}

