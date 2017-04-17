package main

import (
	"github.com/jpg0/flickrup/processing"
	"github.com/jpg0/flickrup/tags"
	"github.com/jpg0/flickrup/archive"
	"github.com/jpg0/flickrup/listen"
	"github.com/jpg0/flickrup/config"
	"github.com/juju/errors"
	flickrupconfig "github.com/jpg0/flickrup/config"
	log "github.com/Sirupsen/logrus"
	"github.com/jpg0/flickrup/flickr"
	"github.com/jpg0/flickrup/filetype"
	"time"
)

const WAIT_TIME = time.Minute * 5

func CreateAndRunPipeline(config *config.Config) error {
	triggerChannel, err := listen.Watch(config)

	if err != nil {
		return errors.Trace(err)
	}

	completions := make(chan struct{})

	l := listen.NewListener(listen.Coalesce(triggerChannel), completions)

	processor, err := ProcessorPipeline(config)
	if err != nil {
		return errors.Trace(err)
	}
	preprocessor, err := PreprocessorPipeline(config)
	if err != nil {
		return errors.Trace(err)
	}

	// initial run
	log.Infof("Triggering initial run...")
	l.TriggerNow()

	for {
		select {
		case beginEvent := <-l.BeginChannel():

			pause := beginEvent.AfterPause

			doRun:
			if pause {
				log.Infof("Waiting for 5 mins...")
				time.Sleep(WAIT_TIME)
			}

			result := SafePerformRun(preprocessor, processor, config)

			if result == RESULT_RERUN {
				log.Infof("Rerunnning")
				pause = false
				goto doRun
			} else if result == RESULT_RESCHEDULE {
				log.Infof("Rescheduling")
				pause = true
				goto doRun
			}

			completions <- struct {}{}
		}
	}

	return nil
}

func ProcessorPipeline(config *flickrupconfig.Config, additionalStages ...processing.Stage) (processing.Processor, error) {

	client, err := flickr.NewUploadClient(config)

	if err != nil {
		return nil, errors.Trace(err)
	}

	tagSetProcessor, err := tags.NewTagSetProcessor(config)
	rewriter := tags.NewRewriter()

	wiredStages := []processing.Stage{
		processing.AsStage(rewriter.MaybeRewrite),
		processing.AsStage(tags.MaybeBlock),
		processing.AsStage(tags.MaybeReplace),
		processing.AsStage(tags.ExtractVisibility),
		filetype.SidecarStage(),
		client.Stage(),
		tagSetProcessor.Stage(),
		processing.AsStage(archive.Archive),
	}

	return processing.Chain(
		append(wiredStages, additionalStages...)...,
	), nil
}

func PreprocessorPipeline(config *flickrupconfig.Config, additionalStages ...processing.PreStage) (processing.Preprocessor, error) {
	return processing.ChainPreStages(
		filetype.VideoConversionStage(),
	), nil
}

func SafePerformRun(preprocessor processing.Preprocessor, processor processing.Processor, config *config.Config) ProcessResult {

	rerun, err := PerformRun(preprocessor, processor, config)

	if err != nil {
		log.Errorf("Run failed: ", err)
	}

	return rerun
}