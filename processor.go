package main

import (
	"github.com/jpg0/flickrup/processing"
	"github.com/jpg0/flickrup/config"
	"io/ioutil"
	"sort"
	"github.com/jpg0/flickrup/filetype"
	"github.com/juju/errors"
	log "github.com/Sirupsen/logrus"
)

func TaggedFileFactory() *processing.TaggedFileFactory {
	return processing.MergeTaggedFileFactories(filetype.TaggedImageFactory(), filetype.TaggedVideoFactory())
}

type ByDateTaken []processing.TaggedFile

func (a ByDateTaken) Len() int           { return len(a) }
func (a ByDateTaken) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByDateTaken) Less(i, j int) bool { return a[i].RealDateTaken().Before(a[j].RealDateTaken()) }

func takeWhile(files []processing.TaggedFile, f func(processing.TaggedFile)bool) []processing.TaggedFile {
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

func PerformRun(processor processing.Processor, config *config.Config) (bool, error) {

	fileInfos, err := ioutil.ReadDir(config.WatchDir)

	if err != nil {
		return false, errors.Trace(err)
	}

	log.Infof("Processing %v files", len(fileInfos))

	//TODO: file conversion

	factory := TaggedFileFactory()
	files := make([]processing.TaggedFile, 0)

	for _, fileInfo := range fileInfos {
		file, err := factory.LoadTaggedFile(config.WatchDir + "/" + fileInfo.Name())

		if err != nil {
			log.Warnf("Failed to load file %v, ignoring", fileInfo.Name())
		} else {
			files = append(files, file)
		}
	}

	log.Infof("Scanned %v files", len(files))

	sort.Sort(ByDateTaken(files))

	byDate := takeWhile(files, func(file processing.TaggedFile) bool {
		return len(file.Keywords().All()) > 0
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
		err = processor(ctx)

		if err != nil {
			log.Warnf("Failed to process %v", toProcess.Name())
			log.Warn(err)
		} else {
			log.Infof("Processing complete for %v", toProcess.Name())
		}
	}

	log.Infof("Processed %v files", len(byDate))

	return false, nil
}
