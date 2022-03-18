package notify

import (
	"context"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/deweppro/go-app/console"
	"github.com/fsnotify/fsnotify"
)

func New(ctx context.Context, dir string, c chan string) {
	watcher, err := fsnotify.NewWatcher()
	console.FatalIfErr(err, "init watcher")
	defer watcher.Close()

	go handler(ctx, watcher, c)

	err = filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if info.Mode().IsDir() && !strings.Contains(path, "/.") {
			return watcher.Add(path)
		}
		return nil
	})
	console.FatalIfErr(err, "add dir [%s] for watch", dir)

	<-ctx.Done()
}

func handler(ctx context.Context, w *fsnotify.Watcher, c chan string) {
	for {
		select {
		case <-ctx.Done():
			return

		case event, ok := <-w.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				c <- event.Name
				console.Infof("%+v", event.Name)
			}

		case err, ok := <-w.Errors:
			if !ok {
				return
			}
			console.Errorf(err.Error())
		}
	}
}
