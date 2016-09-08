package listen

import (
	"github.com/fsnotify/fsnotify"
	"github.com/jpg0/flickrup/config"
	log "github.com/Sirupsen/logrus"
	"github.com/juju/errors"
)

func Watch(cfg *config.Config) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return errors.Trace(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				log.Debug("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Debug("modified file:", event.Name)
					ChangeOccurred()
				}
			case err := <-watcher.Errors:
				log.Error("error:", err)
			}
		}
	}()

	err = watcher.Add(cfg.WatchDir)
	if err != nil {
		return errors.Trace(err)
	}
	<-done

	return nil
}

func ChangeOccurred() {

}