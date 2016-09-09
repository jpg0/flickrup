package listen

import log "github.com/Sirupsen/logrus"

type Listener struct {
	processing bool
	queued     bool
	begin      chan struct{}
}

// Op describes a set of file operations.
type Update uint32

// These are the generalized file operations that can trigger a notification.
const (
	Triggered Update = 1 << iota
	ProcessingComplete
)

func NewListener(triggers <-chan struct{}, completions <-chan struct{}) *Listener {

	begin := make(chan struct{}, 1)

	l := &Listener{begin:begin}

	go func() {
		defer close(l.begin)
		for {
			select {
			case <-completions:
				l.changeOccurred(ProcessingComplete)
			case <-triggers:
				l.changeOccurred(Triggered)
			}
		}
	}()

	return l
}

func (l *Listener) BeginChannel() <-chan struct{} {
	return l.begin
}

func (l *Listener) Trigger() {
	l.changeOccurred(Triggered)
}

//single threaded
func (l *Listener) changeOccurred(u Update) {
	log.Info("Change detected")
	switch u {
	case Triggered:
		if l.processing {
			log.Debug("Processing queued")
			l.queued = true
		} else {
			log.Debug("Processing triggered")
			l.processing = true
			l.begin <- struct{}{}
		}
	case ProcessingComplete:
		if l.processing {
			l.processing = false
			log.Debug("Processing complete")
			if l.queued {
				log.Debug("Queued processing triggered")
				l.queued = false
				l.processing = true
				l.begin <- struct{}{}
			}
		} else {
			log.Errorf("Not marked as processing at completion of processing")
		}
	}
}