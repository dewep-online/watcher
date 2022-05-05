package main

import (
	"github.com/dewep-online/watcher/internal/app"
	"github.com/deweppro/go-app/console"
)

func main() {
	root := console.New("watcher", "help watcher")
	root.AddCommand(app.RunApp())
	root.Exec()
}
