package app

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/dewep-online/watcher/pkg/exc"
	"github.com/dewep-online/watcher/pkg/notify"
	"github.com/deweppro/go-app/application/sys"
	"github.com/deweppro/go-app/console"
)

func RunApp() console.CommandGetter {
	return console.NewCommand(func(setter console.CommandSetter) {
		setter.Setup("run", "executing a command on file changes")
		setter.Example(`run --interval=30 "echo '%file%'"`)
		setter.Argument(1, nil)
		setter.Flag(func(fs console.FlagsSetter) {
			fs.IntVar("interval", 10, "restart interval on changes")
		})
		setter.ExecFunc(func(args []string, ival int64) {
			console.Infof("Run command: %+v", args)

			ctx, cncl := context.WithCancel(context.Background())
			go sys.OnSyscallStop(cncl)

			nn, err := notify.New()
			console.FatalIfErr(err, "init watcher")

			dir, err := os.Getwd()
			console.FatalIfErr(err, "get current dir")
			console.Infof("Dir: %+v", dir)

			err = filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
				if info.Mode().IsDir() && !strings.Contains(path, "/.") {
					return nn.Path(path)
				}
				return nil
			})
			console.FatalIfErr(err, "add dirs")

			nn.Handler(func(lc context.Context) {
				if err := exc.Call(lc, dir, strings.Join(args, " ")); err != nil {
					console.Errorf(err.Error())
				}
			})

			err = nn.Run(ctx, ival)
			console.FatalIfErr(err, "run watcher")
		})
	})
}
