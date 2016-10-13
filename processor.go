package main

import (
	"github.com/jpg0/flickrup/processing"
	"github.com/jpg0/flickrup/config"
	"io/ioutil"
	"sort"
	"github.com/jpg0/flickrup/filetype"
	"github.com/juju/errors"
	log "github.com/Sirupsen/logrus"
	"os"
	"sync"
	"time"
)

func TaggedFileFactory() *processing.TaggedFileFactory {
	return processing.MergeTaggedFileFactories(filetype.TaggedImageFactory(), filetype.TaggedVideoFactory())
}

const FILE_LOAD_CONCURRENCY = 4
var MAX_TIME = time.Unix(2130578470, 0)

type ByDateTaken []processing.TaggedFile

func (a ByDateTaken) Len() int {
	return len(a)
}
func (a ByDateTaken) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a ByDateTaken) Less(i, j int) bool {
	return MustDateTaken(a[i]).Before(MustDateTaken(a[j]))
}

func MustDateTaken(f processing.TaggedFile) time.Time {
	rv := f.DateTaken()

	if rv.IsZero() {
		rv = MAX_TIME
	}

	return rv
}

func takeWhile(files []processing.TaggedFile, f func(processing.TaggedFile) bool) []processing.TaggedFile {
	rv := make([]processing.TaggedFile, 0)

	for _, x := range files {
		if f(x) {
			rv = append(rv, x)
		} else {
			return rv
		}
	}

	return rv
}

func PerformRun(preprocessor processing.Processor, processor processing.Processor, config *config.Config) (bool, error) {

	fileInfos, err := ioutil.ReadDir(config.WatchDir)

	if err != nil {
		return false, errors.Trace(err)
	}

	log.Infof("Processing %v files", len(fileInfos))

	//TODO: file conversion

	files := LoadFiles(fileInfos, TaggedFileFactory(), config)

	log.Infof("Scanned %v files", len(files))

	var result processing.ProcessingResult

	for _, toProcess := range files {
		log.Debugf("Beginning preprocessing for %v", toProcess.Name())
		ctx := processing.NewProcessingContext(config, toProcess)
		result = preprocessor(ctx)

		switch result.ResultType {
		case processing.SuccessResult:
			log.Debugf("Preprocessing complete for %v", toProcess.Name())
		case processing.ErrorResult:
			log.Warnf("Failed to preprocess %v", toProcess.Name())
			log.Warn(result.Error)
		case processing.RestartResult:
			log.Infof("Restarting run after preprocessing %v", toProcess.Name())
			return true, nil
		}
	}

	log.Infof("Preprocessed %v files", len(files))

	sort.Sort(ByDateTaken(files))

	byDate := takeWhile(files, func(file processing.TaggedFile) bool {
		return file.Keywords().All().Size() > 0
	})

	if len(byDate) == 0 {
		log.Info("No files selected for upload")
		if len(files) > 0 {
			log.Infof("Stopped on %v", files[0].Name())
			return false, nil
		}
	} else {
		log.Infof("Selected %v files for upload", len(byDate))
	}

	for _, toProcess := range byDate {
		log.Debugf("Beginning processing for %v", toProcess.Name())
		ctx := processing.NewProcessingContext(config, toProcess)
		result = processor(ctx)

		switch result.ResultType {
		case processing.SuccessResult:
			log.Infof("Processing complete for %v", toProcess.Name())
		case processing.ErrorResult:
			log.Warnf("Failed to process %v", toProcess.Name())
			log.Warn(result.Error)
		case processing.RestartResult:
			log.Infof("Restarting run after processing %v", toProcess.Name())
			return true, nil
		}
	}

	log.Infof("Processed %v files", len(byDate))

	return false, nil
}

func LoadFiles(files []os.FileInfo, factory *processing.TaggedFileFactory, config *config.Config) []processing.TaggedFile {

	type Job struct {
		Index int
		Input os.FileInfo
	}

	processed := make([]processing.TaggedFile, len(files))

	tasks := make(chan Job, 64)

	// spawn N worker goroutines
	var wg sync.WaitGroup
	for i := 0; i < FILE_LOAD_CONCURRENCY; i++ {
		wg.Add(1)
		go func() {
			for job := range tasks {
				file, err := factory.LoadTaggedFile(config.WatchDir + "/" + job.Input.Name())

				if err != nil {
					switch e := err.(type) {
						case processing.NoConstructorAvailableError:
							log.Debugf("Ignoring file %v", job.Input.Name())
						default:
							log.Warnf("Failed to load file %v, ignoring", job.Input.Name())
							log.Warnf(e.Error())
					}

				} else {
					processed[job.Index] = file
				}
			}
			wg.Done()
		}()
	}

	for i, job := range files {
		tasks <- Job{i, job}
	}

	close(tasks)

	// wait for the workers to finish
	wg.Wait()

	rv := make([]processing.TaggedFile, 0)

	for _, completed := range processed {
		if completed != nil {
			rv = append(rv, completed)
		}
	}

	return rv

}


