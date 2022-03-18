package exc

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/deweppro/go-app/console"
)

func Exec(ctx context.Context, cmd, dir string) {
	c := exec.CommandContext(ctx, "/bin/sh", "-xec", fmt.Sprintln(cmd, " <&-"))
	c.Dir = dir
	b, err := c.CombinedOutput()
	console.FatalIfErr(err, string(b))
	fmt.Println(string(b))
}
