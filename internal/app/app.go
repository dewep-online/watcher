package app

import (
	"context"
	"os"

	"github.com/dewep-online/watcher/pkg/notify"
	"github.com/deweppro/go-app/application"
	"github.com/deweppro/go-app/console"
)

func RunApp() console.CommandGetter {
	return console.NewCommand(func(setter console.CommandSetter) {
		setter.Setup("run", "executing a command on file changes")
		setter.Example(`run --interval=30 "echo '%file%'"`)
		setter.Argument(1, nil)
		setter.Flag(func(fs console.FlagsSetter) {
			fs.IntVar("interval", 30, "restart interval on changes")
		})
		setter.ExecFunc(func(args []string, _ int64) {
			console.Infof("Run command: %+v", args)

			dir, err := os.Getwd()
			console.FatalIfErr(err, "get current dir")
			console.Infof("Dir: %+v", dir)

			ctx, cncl := context.WithCancel(context.Background())
			go application.OnSyscallStop(cncl)

			c := make(chan string, 100)

			notify.New(ctx, dir, c)

			<-ctx.Done()
		})
	})
}
