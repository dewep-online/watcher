package app

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync/atomic"
	"time"

	"github.com/deweppro/go-app/application"
	"github.com/deweppro/go-app/console"
	"github.com/fsnotify/fsnotify"
)

func RunApp() console.CommandGetter {
	return console.NewCommand(func(setter console.CommandSetter) {
		setter.Setup("run", "run command")
		setter.Example("run \"echo 'test'\"")
		setter.Argument(1, nil)
		setter.ExecFunc(func(args []string) {
			console.Infof("Run command: %+v", args)

			ctx, cncl := context.WithCancel(context.Background())

			watcher, err := fsnotify.NewWatcher()
			console.FatalIfErr(err, "init watcher")
			defer watcher.Close()

			dir, err := os.Getwd()
			console.FatalIfErr(err, "get current dir")
			console.Infof("Dir: %+v", dir)

			go application.OnSyscallStop(cncl)
			go handler(watcher, ctx, cncl, func() {
				command(strings.Join(args, " "), dir)
			})

			err = filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
				if info.Mode().IsDir() && !strings.HasSuffix(path, ".git") {
					return watcher.Add(path)
				}
				return nil
			})
			console.FatalIfErr(err, "add dir [%s] for watch", dir)

			<-ctx.Done()
		})
	})
}

var lock uint32 = 0

func handler(w *fsnotify.Watcher, ctx context.Context, cncl context.CancelFunc, call func()) {
	defer cncl()

	for {
		select {
		case <-ctx.Done():
			return

		case event := <-w.Events:
			if event.Op&fsnotify.Write == fsnotify.Write {
				console.Infof("%+v", event.Name)
				if atomic.LoadUint32(&lock) == 0 {
					go call()
				}
			}

		case err := <-w.Errors:
			console.Errorf(err.Error())

		}
	}
}

func command(cmd, dir string) {
	atomic.StoreUint32(&lock, 1)
	defer atomic.StoreUint32(&lock, 0)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	c := exec.CommandContext(ctx, "/bin/sh", "-xec", fmt.Sprintln(cmd, " <&-"))
	if dir, err := os.Getwd(); err == nil {
		c.Dir = dir
	}
	b, err := c.CombinedOutput()
	console.FatalIfErr(err, string(b))

	fmt.Println(string(b))
}
