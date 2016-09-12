package tags


import (
	"testing"
	"github.com/jpg0/flickrup/processing"
	"github.com/jpg0/flickrup/config"
	"github.com/golang/mock/gomock"
	"github.com/jpg0/flickrup/mocks"
	"github.com/jpg0/flickrup/testlib"
	"time"
)

func TestNoSetsRequired(t *testing.T){
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockObj := mock_flickr.NewMockSetClient(mockCtrl)

	config := &config.Config{
		TagsetPrefix: "1:2:3=",
	}
	file := testlib.NewFakeTaggedFile("", "", make(map[string]string), nil, time.Time{})
	ctx := processing.NewProcessingContext(config, file)


	tsp := &TagSetProcessor{
		setClient: mockObj,
	}

	//run
	tsp.Stage()(ctx, processing.SuccessProcessor)
}

func TestAddToSetRequired(t *testing.T){

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockObj := mock_flickr.NewMockSetClient(mockCtrl)

	setTime := time.Now()
	fileTime := setTime.Add(time.Hour)

	config := &config.Config{
		TagsetPrefix: "1:2:3=",
	}
	file := testlib.NewFakeTaggedFile("", "", make(map[string]string), []string{"1:2:3=X"}, fileTime)
	ctx := processing.NewProcessingContext(config, file)

	ctx.UploadedId = "test_id"

	mockObj.EXPECT().DateOfSet("X").Times(1).Return(setTime, nil)
	mockObj.EXPECT().AddToSet(ctx.UploadedId, "X", fileTime)

	tsp := &TagSetProcessor{
		setClient: mockObj,
	}

	//run
	tsp.Stage()(ctx, processing.SuccessProcessor)
}