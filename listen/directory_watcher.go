package listen

import (
	"github.com/fsnotify/fsnotify"
	"github.com/jpg0/flickrup/config"
	log "github.com/Sirupsen/logrus"
	"github.com/juju/errors"
)

func Watch(cfg *config.Config) (<-chan struct {}, error){
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, errors.Trace(err)
	}
	c := make(chan struct{})

	go func() {
		for {
			select {
			case <-watcher.Events:
				//log.Debugf("Detected Change:", e)
				c <- struct {}{}
			case err := <-watcher.Errors:
				log.Error("error:", err)
			}
		}
	}()

	err = watcher.Add(cfg.WatchDir)
	if err != nil {
		return nil, errors.Trace(err)
	}

	return c, nil
}