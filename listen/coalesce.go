package listen

import "time"

func Coalesce(in <-chan struct {}) <-chan struct {} {

	out := make(chan struct{})

	timer := time.NewTimer(0)

	var timerCh <-chan time.Time
	var outCh chan<- struct {}

	go func() {
		for {
			select {
			case <-in:
				if timerCh == nil {
					timer.Reset(time.Second)
					timerCh = timer.C
				}
			case <-timerCh:
				outCh = out
				timerCh = nil
			case outCh <- struct{}{}:
				outCh = nil
			}
		}
	}()

	return out
}