package listen

import (
	"github.com/fsnotify/fsnotify"
	"github.com/jpg0/flickrup/config"
	log "github.com/Sirupsen/logrus"
	"github.com/juju/errors"
)

func Watch(cfg *config.Config, cm *ChangeManger) (<-chan struct {}, error){
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, errors.Trace(err)
	}
	c := make(chan struct{})

	go func() {
		for {
			select {
			case e := <-watcher.Events:
				if !cm.ChangeObserved(cfg.WatchDir, e.Name) {
					log.Debugf("Observed file change: %v", e.Name)
					c <- struct{}{}
				} else {
					log.Debugf("Ignoring status file change: %v", e.Name)
				}
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