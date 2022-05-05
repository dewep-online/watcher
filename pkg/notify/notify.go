package notify

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/deweppro/go-app/console"
	"github.com/fsnotify/fsnotify"
)

type Notify struct {
	watcher *fsnotify.Watcher
	handler func(ctx context.Context)
}

func New() (*Notify, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("init watcher: %w", err)
	}
	return &Notify{
		watcher: watcher,
	}, nil
}

func (v *Notify) Path(path string) error {
	return v.watcher.Add(path)
}

func (v *Notify) Handler(cb func(ctx context.Context)) {
	v.handler = cb
}

func (v *Notify) Run(ctx context.Context, ival int64) error {
	tick := time.NewTicker(time.Second * time.Duration(ival))
	defer tick.Stop()

	func() {
		var change int32 = 0

		ctx0, cncl0 := context.WithCancel(ctx)
		go v.handler(ctx0)

		for {
			select {
			case <-tick.C:
				if atomic.CompareAndSwapInt32(&change, 1, 0) {
					cncl0()
					<-ctx0.Done()

					ctx0, cncl0 = context.WithCancel(ctx)
					go v.handler(ctx0)
				}

			case <-ctx.Done():
				cncl0()
				return

			case ev, ok := <-v.watcher.Events:
				if !ok {
					cncl0()
					return
				}
				if ev.Op&fsnotify.Write == fsnotify.Write {
					atomic.CompareAndSwapInt32(&change, 0, 1)
				}

			case ev, ok := <-v.watcher.Errors:
				if !ok {
					cncl0()
					return
				}
				console.Errorf(ev.Error())
			}
		}
	}()

	return v.watcher.Close()
}
